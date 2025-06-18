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

/* A micro-Kanren in go, based on the following implementation found in [1]
in Scheme, by the original authors:

Variables are represented as vectors (we'll represent them as refs)

(define (var c) (vector c))
(define (var? x) (vector? x))
(define (var=? x1 x2)
  (= (vector-ref x1 0) (vector-ref x2 0)))

There is no check for circularities when searching for a variable's value in a substitution

(define (walk u s)
	(let ((pr (and (var? u) (assp (λ (v) (var=? u v)) s))))
		(if pr (walk (cdr pr) s) u)))

(define (ext-s x v s)
	`((,x . ,v) . ,s))

The first of the goal constructors, ≡, returns a goal if the two arguments unify.

(define (≡ u v)
	(λ_g (s/c)
		(let ((s (unify u v (car s/c))))
			(if s (unit `(,s . ,(cdr s/c))) mzero))))

(define (unit s/c) (cons s/c mzero))
(define mzero '())

Basic unification, the core of any logic programming language.

(define (unify u v s)
	(let ((u (walk u s)) (v (walk v s)))
		(cond
			((and (var? u) (var? v) (var=? u v)) s)
			((var? u) (ext- s u v s))
			((var? v) (ext- s v u s))
			((and (pair? u) (pair? v))
			 (let ((s (unify (car u) (car v) s)))
				 (and s (unify (cdr u) (cdr v) s))))
			(else (and (eqv? u v) s)))))

Creates a new variable, and is also a goal constructor.

(define (call/fresh f)
	(λ_g (s/c)
		(let ((c (cdr s/c)))
			((f (var c)) `(,(car s/c) . ,(+ c 1))))))

Disjunction and Conjunction are goal constructors that combine goals.

(define (disj g1 g2)
	(λ_g (s/c) (mplus (g1 s/c) (g2 s/c))))
(define (conj g1 g2)
	(λ_g (s/c) (bind (g1 s/c) g2)))

Merges two streams

(define (mplus $1 $2)
	(cond
		((null? $1) $2)
		((procedure? $1) (λ_$ () (mplus $2 ($1))))
		(else (cons (car $1) (mplus (cdr $1) $2)))))

Used in `conj` to apply a goal from one stream to all goals in a second stream.

(define (bind $ g)
	(cond
		((null? $) mzero)
		((procedure? $) (λ_$ () (bind ($) g)))
		(else (mplus (g (car $)) (bind (cdr $) g)))))


[1]: Hemann, J., & Friedman, D. P. (2013). "microkanren: A Minimal Functional
     Core for Relational Programming". Workshop on Scheme and Functional
		 Programming, Alexandria, United States.
*/
