// Copyright (c) 2025 Kevin Damm
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package ukanren

// A micro-Kanren in go, based on the implementation found in:
//   Hemann, J., & Friedman, D. P. (2013). "microkanren: A Minimal Functional
//   Core for Relational Programming". Workshop on Scheme and Functional
//   Programming, Alexandria, United States.

// All representable values are atoms.
type Atom interface {
	isAtom()
}

func isVar(v Atom) bool {
	switch v.(type) {
	case Variable:
		return true
	default:
		return false
	}
}

// The only literal type in this implementation is Integer (more can be added).
type Literal int

func (lit Literal) isAtom() {}

// Variable representation, using de Bruijn indices.
type Variable struct {
	value int
}

func (v Variable) isAtom() {}

// List constructor.
type Cons []Atom

func (Cons) isAtom() {}

func (list Cons) car() Atom {
	if len(list) > 0 {
		return list[0]
	}
	return nil
}

func (list Cons) cdr() Atom {
	if len(list) > 1 {
		return list[1:]
	}
	return nil
}

// Variable substitutions are represented by a mapping of atoms to atoms.
type Subs map[Atom]Atom

func (subs Subs) extend(key Atom, val Atom) Subs {
	copy := make(Subs)
	for k, v := range subs {
		copy[k] = v
	}
	copy[key] = val
	return copy
}

func (subs Subs) intend(term Atom) Atom {
	if isVar(term) && subs[term] != nil {
		return subs.intend(subs[term])
	}
	return term
}

// Unification, the core of any logical language.  Add to this if adding types.
func unify(t1 Atom, t2 Atom, subs Subs) (Subs, bool) {
	t1 = subs.intend(t1)
	t2 = subs.intend(t2)

	// Simple case of identity, return the same mapping.
	if isVar(t1) && isVar(t2) && t1 == t2 {
		return subs, true
	}

	// If t1 is a variable, unify with an updated mapping, t1 adopting t2.
	if isVar(t1) {
		return subs.extend(t1, t2), true
	}

	// If t1 is a variable, unify with an updated mapping, t2 adopting t1.
	if isVar(t2) {
		return subs.extend(t2, t1), true
	}

	// If both are lists, try to unify each of their elements.
	list1, ok1 := t1.(Cons)
	list2, ok2 := t2.(Cons)
	if ok1 && ok2 {
		updated, ok := unify(list1.car(), list2.car(), subs)
		if ok {
			return unify(list1.cdr(), list2.cdr(), updated)
		}
		return nil, false
	}

	// If they are not variables or lists, but are equal, they already unify.
	if t1 == t2 {
		return subs, true
	}

	// We've run out of ways it could successfully unify; return failure.
	return nil, false
}

// Each State is a successful resolution, returned as elements of a Stream.
type State struct {
	vars  Subs
	count int
}

func (s State) isAtom() {}

// Streams are the return value of conjunction and disjunction.
type Stream Cons

// (internal) constructor for streams as a result of finding and combining goals.
func newStream(state State) Stream {
	stream := make([]Atom, 1)
	stream[0] = state
	return stream
}

// Goals are closures for converting a state into a stream.
type Goal func(State) Stream

func EvalGoal(goal Goal) Stream {
	state := State{make(Subs), 0}
	return goal(state)
}

func EvalFresh(fn func(Atom) Goal) Goal {
	return func(state State) Stream {
		goal := fn(Variable{state.count})
		nextState := State{state.vars, state.count + 1}
		return goal(nextState)
	}
}

// Combine all elements of two streams together.
func append(s1 Stream, s2 Stream) Stream {
	if s1[0] == nil {
		return s2
	}
	stream := make([]Atom, len(s1)+len(s2))
	copy(stream, s1)
	copy(stream[len(s1):], s2)
	return stream
}

// Combine all results of a goal if they exist in the provided stream.
func mappend(goal Goal, s Stream) Stream {
	if len(s) == 0 {
		// Nothing to join with.
		return nil
	}

	h := s[0]
	state, ok := h.(State)
	if !ok {
		return nil
	}

	stream := append(goal(state), mappend(goal, s[1:]))
	return stream
}

// Equality is tested by term unification.
func Equal(t1 Atom, t2 Atom) Goal {
	return func(state State) Stream {
		subs, ok := unify(t1, t2, state.vars)
		if ok {
			return newStream(State{subs, state.count})
		}
		return nil
	}
}

// Disjunction function, accept all results of either goal.
func Disj(g1 Goal, g2 Goal) Goal {
	return func(state State) Stream {
		return append(g1(state), g2(state))
	}
}

// Conjunction function, accept only states which are results of both goals.
func Conj(g1 Goal, g2 Goal) Goal {
	return func(state State) Stream {
		return mappend(g2, g1(state))
	}
}
