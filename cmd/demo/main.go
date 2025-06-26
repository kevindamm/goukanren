// Copyright (c) 2025 Kevin Damm.
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

func main() {
	// A classic logical problem, the {Chicken, Goat, Cabbage} puzzle, formulated
	// as a search query for valid paths which satisfy the following constraints:
	//
	// ```A farmer with a wolf, a goat, and a cabbage must cross a river by boat.
	// The boat can carry only the farmer and a single item. If left unattended
	// together, the wolf would eat the goat, or the goat would eat the cabbage.
	// How can they cross the river without anything being eaten?```

	// In miniKanren this would be the following:

	// (defrel (participant/o x)
	//   (member/o x '(Goat Wolf Cabbage)))

	// (defrel (eats/o who what)
	//   (cond/e
	//     ((== who 'Goat) (== what 'Cabbage))
	//     ((== who 'Wolf) (== what 'Goat))))

	// Represent the group on one side of the river as a set-like relation group.
	//
	// (defrel (group-helper/o g seen)
	//   (cond/e
	//     ((== g '()))
	//     ((fresh (head tail)
	//        (cons/o head tail g)
	//        (participant/o head)
	//        (cond/a
	//          ((once (member/o head seen)) fail)
	//          ((fresh (updated)
	//             (cons/o head seen updated)
	//             (group-helper/o tail updated))))))))
	//
	// Include a base case for the constructive definition above.
	//
	// (defrel (group/o g)
	//   (group-helper/o g '()))

	// Safe groups
	//
	// (defrel (safe-group/o g)
	//   (group/o g)
	//   (cond/u ((fresh (who what)
	//              (eats/o who what)
	//              (member/o who g)
	//              (member/o what g))
	//            ;; fail if we find a pair of participants who eat each other...
	//            fail)
	//           (succeed)))

	// The selection process
	//
	// (defrel (pick/o g participant out)
	//   (group/o g)
	//   (fresh (head tail)
	//     (cons/o head tail g)
	//     (cond/e
	//       ((== participant head)
	//        (== out tail))
	//       ((fresh (intermediate)
	//          (pick/o tail participant intermediate)
	//          (cons/o head intermediate out)))))
	//   (safe-group/o out))

	// The key to solving is knowing that the boat can ferry nothing across the river.
	//
	// (defrel (pick-maybe-nothing/o g participant out)
	//  (cond/e
	//    ((safe-group/o g)
	//     (== participant 'Nothing)
	//     (== out g))
	//    ((pick/o g participant out))))

	// Avoid looping indefinitely by keeping track of groups we've already seen.
	//
	// (defrel (seen-group/o g history)
	//  (fresh (head tail)
	//    (cons/o head tail history)
	//    (cond/e ((same/o g head))
	//            ((seen-group/o g tail)))))

	// (defrel (goal-state/o state)
	//   (fresh (tag group goal history)
	//     (== `(River-Bank ,tag ,group ,goal ,history) state)
	//     (same/o group goal)))

	// (defrel (do-plan/o from to plan)
	//	 (cond/e
	//	  ((goal-state/o from)
	//	   (== plan '()))
	//	  ((fresh (participant step remaining updated-from updated-to)
	//	     (pick-participant-from/o from participant updated-from)
	//	     (add-participant-to/o to participant updated-to)
	//	     (make-step/o from to participant step)
	//	     (cons/o step remaining plan)
	//	     (do-plan/o updated-to updated-from remaining)))))
	//
	// (defrel (wolf-goat-cabbage-plan/o plan)
	//   (let ([from '(Wolf Goat Cabbage)]
	//         [to '()])
	//     (do-plan/o
	//      `(River-Bank Left ,from ,to (,from))
	//      `(River-Bank Right ,to ,from (,to))
	//      plan)))
}
