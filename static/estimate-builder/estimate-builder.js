//#region node_modules/svelte/src/internal/shared/utils.js
var e = Array.isArray, t = Array.prototype.indexOf, n = Array.prototype.includes, r = Array.from;
Object.keys;
var i = Object.defineProperty, a = Object.getOwnPropertyDescriptor, o = Object.getOwnPropertyDescriptors, s = Object.prototype, c = Array.prototype, l = Object.getPrototypeOf, u = Object.isExtensible, d = () => {};
function f(e) {
	for (var t = 0; t < e.length; t++) e[t]();
}
function p() {
	var e, t;
	return {
		promise: new Promise((n, r) => {
			e = n, t = r;
		}),
		resolve: e,
		reject: t
	};
}
function m(e, t) {
	if (Array.isArray(e)) return e;
	if (t === void 0 || !(Symbol.iterator in e)) return Array.from(e);
	let n = [];
	for (let r of e) if (n.push(r), n.length === t) break;
	return n;
}
var h = 1024, g = 2048, _ = 4096, v = 8192, y = 16384, b = 32768, x = 1 << 25, S = 65536, C = 1 << 19, w = 1 << 20, T = 1 << 25, E = 65536, D = 1 << 21, ee = 1 << 22, te = 1 << 23, ne = Symbol("$state"), re = Symbol("legacy props"), ie = Symbol(""), ae = new class extends Error {
	name = "StaleReactionError";
	message = "The reaction that called `getAbortSignal()` was re-run or destroyed";
}(), oe = !!globalThis.document?.contentType && /* @__PURE__ */ globalThis.document.contentType.includes("xml");
//#endregion
//#region node_modules/svelte/src/internal/client/errors.js
function se() {
	throw Error("https://svelte.dev/e/async_derived_orphan");
}
function ce(e, t, n) {
	throw Error("https://svelte.dev/e/each_key_duplicate");
}
function le(e) {
	throw Error("https://svelte.dev/e/effect_in_teardown");
}
function ue() {
	throw Error("https://svelte.dev/e/effect_in_unowned_derived");
}
function de(e) {
	throw Error("https://svelte.dev/e/effect_orphan");
}
function fe() {
	throw Error("https://svelte.dev/e/effect_update_depth_exceeded");
}
function pe(e) {
	throw Error("https://svelte.dev/e/props_invalid_value");
}
function me() {
	throw Error("https://svelte.dev/e/state_descriptors_fixed");
}
function he() {
	throw Error("https://svelte.dev/e/state_prototype_fixed");
}
function ge() {
	throw Error("https://svelte.dev/e/state_unsafe_mutation");
}
function _e() {
	throw Error("https://svelte.dev/e/svelte_boundary_reset_onerror");
}
//#endregion
//#region node_modules/svelte/src/constants.js
var ve = {}, O = Symbol(), ye = "http://www.w3.org/1999/xhtml";
function be(e) {
	console.warn("https://svelte.dev/e/hydration_mismatch");
}
function xe() {
	console.warn("https://svelte.dev/e/select_multiple_invalid_value");
}
function Se() {
	console.warn("https://svelte.dev/e/svelte_boundary_reset_noop");
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/hydration.js
var k = !1;
function Ce(e) {
	k = e;
}
var A;
function j(e) {
	if (e === null) throw be(), ve;
	return A = e;
}
function we() {
	return j(/* @__PURE__ */ Zt(A));
}
function M(e) {
	if (k) {
		if (/* @__PURE__ */ Zt(A) !== null) throw be(), ve;
		A = e;
	}
}
function Te(e = 1) {
	if (k) {
		for (var t = e, n = A; t--;) n = /* @__PURE__ */ Zt(n);
		A = n;
	}
}
function Ee(e = !0) {
	for (var t = 0, n = A;;) {
		if (n.nodeType === 8) {
			var r = n.data;
			if (r === "]") {
				if (t === 0) return n;
				--t;
			} else (r === "[" || r === "[!" || r[0] === "[" && !isNaN(Number(r.slice(1)))) && (t += 1);
		}
		var i = /* @__PURE__ */ Zt(n);
		e && n.remove(), n = i;
	}
}
function De(e) {
	if (!e || e.nodeType !== 8) throw be(), ve;
	return e.data;
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/equality.js
function Oe(e) {
	return e === this.v;
}
function ke(e, t) {
	return e == e ? e !== t || typeof e == "object" && !!e || typeof e == "function" : t == t;
}
function Ae(e) {
	return !ke(e, this.v);
}
//#endregion
//#region node_modules/svelte/src/internal/flags/index.js
var je = !1, Me = !1, N = null;
function Ne(e) {
	N = e;
}
function Pe(e, t = !1, n) {
	N = {
		p: N,
		i: !1,
		c: null,
		e: null,
		s: e,
		x: null,
		r: W,
		l: Me && !t ? {
			s: null,
			u: null,
			$: []
		} : null
	};
}
function Fe(e) {
	var t = N, n = t.e;
	if (n !== null) {
		t.e = null;
		for (var r of n) hn(r);
	}
	return e !== void 0 && (t.x = e), t.i = !0, N = t.p, e ?? {};
}
function Ie() {
	return !Me || N !== null && N.l === null;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/task.js
var Le = [];
function Re() {
	var e = Le;
	Le = [], f(e);
}
function ze(e) {
	if (Le.length === 0 && !et) {
		var t = Le;
		queueMicrotask(() => {
			t === Le && Re();
		});
	}
	Le.push(e);
}
function Be() {
	for (; Le.length > 0;) Re();
}
function Ve(e) {
	var t = W;
	if (t === null) return U.f |= te, e;
	if (!(t.f & 32768) && !(t.f & 4)) throw e;
	He(e, t);
}
function He(e, t) {
	for (; t !== null;) {
		if (t.f & 128) {
			if (!(t.f & 32768)) throw e;
			try {
				t.b.error(e);
				return;
			} catch (t) {
				e = t;
			}
		}
		t = t.parent;
	}
	throw e;
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/status.js
var Ue = ~(g | _ | h);
function P(e, t) {
	e.f = e.f & Ue | t;
}
function We(e) {
	e.f & 512 || e.deps === null ? P(e, h) : P(e, _);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/utils.js
function Ge(e) {
	if (e !== null) for (let t of e) !(t.f & 2) || !(t.f & 65536) || (t.f ^= E, Ge(t.deps));
}
function Ke(e, t, n) {
	e.f & 2048 ? t.add(e) : e.f & 4096 && n.add(e), Ge(e.deps), P(e, h);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/store.js
var qe = !1, Je = !1;
function Ye(e) {
	var t = Je;
	try {
		return Je = !1, [e(), Je];
	} finally {
		Je = t;
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/batch.js
var Xe = /* @__PURE__ */ new Set(), F = null, Ze = null, Qe = null, $e = null, et = !1, tt = !1, nt = null, rt = null, it = 0, at = 1, ot = class e {
	id = at++;
	current = /* @__PURE__ */ new Map();
	previous = /* @__PURE__ */ new Map();
	#e = /* @__PURE__ */ new Set();
	#t = /* @__PURE__ */ new Set();
	#n = 0;
	#r = 0;
	#i = null;
	#a = [];
	#o = /* @__PURE__ */ new Set();
	#s = /* @__PURE__ */ new Set();
	#c = /* @__PURE__ */ new Map();
	is_fork = !1;
	#l = !1;
	#u() {
		return this.is_fork || this.#r > 0;
	}
	skip_effect(e) {
		this.#c.has(e) || this.#c.set(e, {
			d: [],
			m: []
		});
	}
	unskip_effect(e) {
		var t = this.#c.get(e);
		if (t) {
			this.#c.delete(e);
			for (var n of t.d) P(n, g), this.schedule(n);
			for (n of t.m) P(n, _), this.schedule(n);
		}
	}
	#d() {
		it++ > 1e3 && ct();
		let t = this.#a;
		this.#a = [], this.apply();
		var n = nt = [], r = [], i = rt = [];
		for (let e of t) try {
			this.#f(e, n, r);
		} catch (t) {
			throw ht(e), t;
		}
		if (F = null, i.length > 0) {
			var a = e.ensure();
			for (let e of i) a.schedule(e);
		}
		if (nt = null, rt = null, this.#u()) {
			this.#p(r), this.#p(n);
			for (let [e, t] of this.#c) mt(e, t);
		} else {
			this.#n === 0 && Xe.delete(this), this.#o.clear(), this.#s.clear();
			for (let e of this.#e) e(this);
			this.#e.clear(), Ze = this, ut(r), ut(n), Ze = null, this.#i?.resolve();
		}
		var o = F;
		if (this.#a.length > 0) {
			let e = o ??= this;
			e.#a.push(...this.#a.filter((t) => !e.#a.includes(t)));
		}
		o !== null && (Xe.add(o), o.#d()), Xe.has(this) || this.#m();
	}
	#f(e, t, n) {
		e.f ^= h;
		for (var r = e.first; r !== null;) {
			var i = r.f, a = (i & 96) != 0;
			if (!(a && i & 1024 || i & 8192 || this.#c.has(r)) && r.fn !== null) {
				a ? r.f ^= h : i & 4 ? t.push(r) : je && i & 16777224 ? n.push(r) : Zn(r) && (i & 16 && this.#s.add(r), nr(r));
				var o = r.first;
				if (o !== null) {
					r = o;
					continue;
				}
			}
			for (; r !== null;) {
				var s = r.next;
				if (s !== null) {
					r = s;
					break;
				}
				r = r.parent;
			}
		}
	}
	#p(e) {
		for (var t = 0; t < e.length; t += 1) Ke(e[t], this.#o, this.#s);
	}
	capture(e, t) {
		t !== O && !this.previous.has(e) && this.previous.set(e, t), e.f & 8388608 || (this.current.set(e, e.v), Qe?.set(e, e.v));
	}
	activate() {
		F = this;
	}
	deactivate() {
		F = null, Qe = null;
	}
	flush() {
		try {
			if (tt = !0, F = this, !this.#u()) {
				for (let e of this.#o) this.#s.delete(e), P(e, g), this.schedule(e);
				for (let e of this.#s) P(e, _), this.schedule(e);
			}
			this.#d();
		} finally {
			it = 0, $e = null, nt = null, rt = null, tt = !1, F = null, Qe = null, Pt.clear();
		}
	}
	discard() {
		for (let e of this.#t) e(this);
		this.#t.clear();
	}
	#m() {
		for (let s of Xe) {
			var e = s.id < this.id, t = [];
			for (let [n, r] of this.current) {
				if (s.current.has(n)) if (e && r !== s.current.get(n)) s.current.set(n, r);
				else continue;
				t.push(n);
			}
			if (t.length !== 0) {
				var n = [...s.current.keys()].filter((e) => !this.current.has(e));
				if (n.length > 0) {
					s.activate();
					var r = /* @__PURE__ */ new Set(), i = /* @__PURE__ */ new Map();
					for (var a of t) dt(a, n, r, i);
					if (s.#a.length > 0) {
						s.apply();
						for (var o of s.#a) s.#f(o, [], []);
					}
					s.deactivate();
				}
			}
		}
	}
	increment(e) {
		this.#n += 1, e && (this.#r += 1);
	}
	decrement(e, t) {
		--this.#n, e && --this.#r, !(this.#l || t) && (this.#l = !0, ze(() => {
			this.#l = !1, this.flush();
		}));
	}
	oncommit(e) {
		this.#e.add(e);
	}
	ondiscard(e) {
		this.#t.add(e);
	}
	settled() {
		return (this.#i ??= p()).promise;
	}
	static ensure() {
		if (F === null) {
			let t = F = new e();
			tt || (Xe.add(F), et || ze(() => {
				F === t && t.flush();
			}));
		}
		return F;
	}
	apply() {
		if (!je || !this.is_fork && Xe.size === 1) {
			Qe = null;
			return;
		}
		Qe = new Map(this.current);
		for (let e of Xe) if (e !== this) for (let [t, n] of e.previous) Qe.has(t) || Qe.set(t, n);
	}
	schedule(e) {
		if ($e = e, e.b?.is_pending && e.f & 16777228 && !(e.f & 32768)) {
			e.b.defer_effect(e);
			return;
		}
		for (var t = e; t.parent !== null;) {
			t = t.parent;
			var n = t.f;
			if (nt !== null && t === W && (je || (U === null || !(U.f & 2)) && !qe)) return;
			if (n & 96) {
				if (!(n & 1024)) return;
				t.f ^= h;
			}
		}
		this.#a.push(t);
	}
};
function st(e) {
	var t = et;
	et = !0;
	try {
		var n;
		for (e && (F !== null && !F.is_fork && F.flush(), n = e());;) {
			if (Be(), F === null) return n;
			F.flush();
		}
	} finally {
		et = t;
	}
}
function ct() {
	try {
		fe();
	} catch (e) {
		He(e, $e);
	}
}
var lt = null;
function ut(e) {
	var t = e.length;
	if (t !== 0) {
		for (var n = 0; n < t;) {
			var r = e[n++];
			if (!(r.f & 24576) && Zn(r) && (lt = /* @__PURE__ */ new Set(), nr(r), r.deps === null && r.first === null && r.nodes === null && r.teardown === null && r.ac === null && Dn(r), lt?.size > 0)) {
				Pt.clear();
				for (let e of lt) {
					if (e.f & 24576) continue;
					let t = [e], n = e.parent;
					for (; n !== null;) lt.has(n) && (lt.delete(n), t.push(n)), n = n.parent;
					for (let e = t.length - 1; e >= 0; e--) {
						let n = t[e];
						n.f & 24576 || nr(n);
					}
				}
				lt.clear();
			}
		}
		lt = null;
	}
}
function dt(e, t, n, r) {
	if (!n.has(e) && (n.add(e), e.reactions !== null)) for (let i of e.reactions) {
		let e = i.f;
		e & 2 ? dt(i, t, n, r) : e & 4194320 && !(e & 2048) && ft(i, t, r) && (P(i, g), pt(i));
	}
}
function ft(e, t, r) {
	let i = r.get(e);
	if (i !== void 0) return i;
	if (e.deps !== null) for (let i of e.deps) {
		if (n.call(t, i)) return !0;
		if (i.f & 2 && ft(i, t, r)) return r.set(i, !0), !0;
	}
	return r.set(e, !1), !1;
}
function pt(e) {
	F.schedule(e);
}
function mt(e, t) {
	if (!(e.f & 32 && e.f & 1024)) {
		e.f & 2048 ? t.d.push(e) : e.f & 4096 && t.m.push(e), P(e, h);
		for (var n = e.first; n !== null;) mt(n, t), n = n.next;
	}
}
function ht(e) {
	P(e, h);
	for (var t = e.first; t !== null;) ht(t), t = t.next;
}
//#endregion
//#region node_modules/svelte/src/reactivity/create-subscriber.js
function gt(e) {
	let t = 0, n = It(0), r;
	return () => {
		fn() && (G(n), yn(() => (t === 0 && (r = or(() => e(() => Bt(n)))), t += 1, () => {
			ze(() => {
				--t, t === 0 && (r?.(), r = void 0, Bt(n));
			});
		})));
	};
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/boundary.js
var _t = S | C;
function vt(e, t, n, r) {
	new yt(e, t, n, r);
}
var yt = class {
	parent;
	is_pending = !1;
	transform_error;
	#e;
	#t = k ? A : null;
	#n;
	#r;
	#i;
	#a = null;
	#o = null;
	#s = null;
	#c = null;
	#l = 0;
	#u = 0;
	#d = !1;
	#f = /* @__PURE__ */ new Set();
	#p = /* @__PURE__ */ new Set();
	#m = null;
	#h = gt(() => (this.#m = It(this.#l), () => {
		this.#m = null;
	}));
	constructor(e, t, n, r) {
		this.#e = e, this.#n = t, this.#r = (e) => {
			var t = W;
			t.b = this, t.f |= 128, n(e);
		}, this.parent = W.b, this.transform_error = r ?? this.parent?.transform_error ?? ((e) => e), this.#i = bn(() => {
			if (k) {
				let e = this.#t;
				we();
				let t = e.data === "[!";
				if (e.data.startsWith("[?")) {
					let t = JSON.parse(e.data.slice(2));
					this.#_(t);
				} else t ? this.#v() : this.#g();
			} else this.#y();
		}, _t), k && (this.#e = A);
	}
	#g() {
		try {
			this.#a = xn(() => this.#r(this.#e));
		} catch (e) {
			this.error(e);
		}
	}
	#_(e) {
		let t = this.#n.failed;
		t && (this.#s = xn(() => {
			t(this.#e, () => e, () => () => {});
		}));
	}
	#v() {
		let e = this.#n.pending;
		e && (this.is_pending = !0, this.#o = xn(() => e(this.#e)), ze(() => {
			var e = this.#c = document.createDocumentFragment(), t = Yt();
			e.append(t), this.#a = this.#x(() => xn(() => this.#r(t))), this.#u === 0 && (this.#e.before(e), this.#c = null, On(this.#o, () => {
				this.#o = null;
			}), this.#b(F));
		}));
	}
	#y() {
		try {
			if (this.is_pending = this.has_pending_snippet(), this.#u = 0, this.#l = 0, this.#a = xn(() => {
				this.#r(this.#e);
			}), this.#u > 0) {
				var e = this.#c = document.createDocumentFragment();
				Mn(this.#a, e);
				let t = this.#n.pending;
				this.#o = xn(() => t(this.#e));
			} else this.#b(F);
		} catch (e) {
			this.error(e);
		}
	}
	#b(e) {
		this.is_pending = !1;
		for (let t of this.#f) P(t, g), e.schedule(t);
		for (let t of this.#p) P(t, _), e.schedule(t);
		this.#f.clear(), this.#p.clear();
	}
	defer_effect(e) {
		Ke(e, this.#f, this.#p);
	}
	is_rendered() {
		return !this.is_pending && (!this.parent || this.parent.is_rendered());
	}
	has_pending_snippet() {
		return !!this.#n.pending;
	}
	#x(e) {
		var t = W, n = U, r = N;
		zn(this.#i), Rn(this.#i), Ne(this.#i.ctx);
		try {
			return ot.ensure(), e();
		} catch (e) {
			return Ve(e), null;
		} finally {
			zn(t), Rn(n), Ne(r);
		}
	}
	#S(e, t) {
		if (!this.has_pending_snippet()) {
			this.parent && this.parent.#S(e, t);
			return;
		}
		this.#u += e, this.#u === 0 && (this.#b(t), this.#o && On(this.#o, () => {
			this.#o = null;
		}), this.#c &&= (this.#e.before(this.#c), null));
	}
	update_pending_count(e, t) {
		this.#S(e, t), this.#l += e, !(!this.#m || this.#d) && (this.#d = !0, ze(() => {
			this.#d = !1, this.#m && Rt(this.#m, this.#l);
		}));
	}
	get_effect_pending() {
		return this.#h(), G(this.#m);
	}
	error(e) {
		var t = this.#n.onerror;
		let n = this.#n.failed;
		if (!t && !n) throw e;
		this.#a &&= (Tn(this.#a), null), this.#o &&= (Tn(this.#o), null), this.#s &&= (Tn(this.#s), null), k && (j(this.#t), Te(), j(Ee()));
		var r = !1, i = !1;
		let a = () => {
			if (r) {
				Se();
				return;
			}
			r = !0, i && _e(), this.#s !== null && On(this.#s, () => {
				this.#s = null;
			}), this.#x(() => {
				this.#y();
			});
		}, o = (e) => {
			try {
				i = !0, t?.(e, a), i = !1;
			} catch (e) {
				He(e, this.#i && this.#i.parent);
			}
			n && (this.#s = this.#x(() => {
				try {
					return xn(() => {
						var t = W;
						t.b = this, t.f |= 128, n(this.#e, () => e, () => a);
					});
				} catch (e) {
					return He(e, this.#i.parent), null;
				}
			}));
		};
		ze(() => {
			var t;
			try {
				t = this.transform_error(e);
			} catch (e) {
				He(e, this.#i && this.#i.parent);
				return;
			}
			typeof t == "object" && t && typeof t.then == "function" ? t.then(o, (e) => He(e, this.#i && this.#i.parent)) : o(t);
		});
	}
};
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/async.js
function bt(e, t, n, r) {
	let i = Ie() ? wt : Et;
	var a = e.filter((e) => !e.settled);
	if (n.length === 0 && a.length === 0) {
		r(t.map(i));
		return;
	}
	var o = W, s = xt(), c = a.length === 1 ? a[0].promise : a.length > 1 ? Promise.all(a.map((e) => e.promise)) : null;
	function l(e) {
		s();
		try {
			r(e);
		} catch (e) {
			o.f & 16384 || He(e, o);
		}
		St();
	}
	if (n.length === 0) {
		c.then(() => l(t.map(i)));
		return;
	}
	var u = Ct();
	function d() {
		Promise.all(n.map((e) => /* @__PURE__ */ Tt(e))).then((e) => l([...t.map(i), ...e])).catch((e) => He(e, o)).finally(() => u());
	}
	c ? c.then(() => {
		s(), d(), St();
	}) : d();
}
function xt() {
	var e = W, t = U, n = N, r = F;
	return function(i = !0) {
		zn(e), Rn(t), Ne(n), i && !(e.f & 16384) && (r?.activate(), r?.apply());
	};
}
function St(e = !0) {
	zn(null), Rn(null), Ne(null), e && F?.deactivate();
}
function Ct() {
	var e = W.b, t = F, n = e.is_rendered();
	return e.update_pending_count(1, t), t.increment(n), (r = !1) => {
		e.update_pending_count(-1, t), t.decrement(n, r);
	};
}
/* @__NO_SIDE_EFFECTS__ */
function wt(e) {
	var t = 2 | g, n = U !== null && U.f & 2 ? U : null;
	return W !== null && (W.f |= C), {
		ctx: N,
		deps: null,
		effects: null,
		equals: Oe,
		f: t,
		fn: e,
		reactions: null,
		rv: 0,
		v: O,
		wv: 0,
		parent: n ?? W,
		ac: null
	};
}
/* @__NO_SIDE_EFFECTS__ */
function Tt(e, t, n) {
	let r = W;
	r === null && se();
	var i = void 0, a = It(O), o = !U, s = /* @__PURE__ */ new Map();
	return vn(() => {
		var t = W, n = p();
		i = n.promise;
		try {
			Promise.resolve(e()).then(n.resolve, n.reject).finally(St);
		} catch (e) {
			n.reject(e), St();
		}
		var c = F;
		if (o) {
			if (t.f & 32768) var l = Ct();
			if (r.b.is_rendered()) s.get(c)?.reject(ae), s.delete(c);
			else {
				for (let e of s.values()) e.reject(ae);
				s.clear();
			}
			s.set(c, n);
		}
		let u = (e, n = void 0) => {
			if (l && l(n === ae), !(n === ae || t.f & 16384)) {
				if (c.activate(), n) a.f |= te, Rt(a, n);
				else {
					a.f & 8388608 && (a.f ^= te), Rt(a, e);
					for (let [e, t] of s) {
						if (s.delete(e), e === c) break;
						t.reject(ae);
					}
				}
				c.deactivate();
			}
		};
		n.promise.then(u, (e) => u(null, e || "unknown"));
	}), pn(() => {
		for (let e of s.values()) e.reject(ae);
	}), new Promise((e) => {
		function t(n) {
			function r() {
				n === i ? e(a) : t(i);
			}
			n.then(r, r);
		}
		t(i);
	});
}
/* @__NO_SIDE_EFFECTS__ */
function I(e) {
	let t = /* @__PURE__ */ wt(e);
	return je || Vn(t), t;
}
/* @__NO_SIDE_EFFECTS__ */
function Et(e) {
	let t = /* @__PURE__ */ wt(e);
	return t.equals = Ae, t;
}
function Dt(e) {
	var t = e.effects;
	if (t !== null) {
		e.effects = null;
		for (var n = 0; n < t.length; n += 1) Tn(t[n]);
	}
}
function Ot(e) {
	for (var t = e.parent; t !== null;) {
		if (!(t.f & 2)) return t.f & 16384 ? null : t;
		t = t.parent;
	}
	return null;
}
function kt(e) {
	var t, n = W;
	zn(Ot(e));
	try {
		e.f &= ~E, Dt(e), t = $n(e);
	} finally {
		zn(n);
	}
	return t;
}
function At(e) {
	var t = kt(e);
	if (!e.equals(t) && (e.wv = Xn(), (!F?.is_fork || e.deps === null) && (e.v = t, e.deps === null))) {
		P(e, h);
		return;
	}
	Fn || (Qe === null ? We(e) : (fn() || F?.is_fork) && Qe.set(e, t));
}
function jt(e) {
	if (e.effects !== null) for (let t of e.effects) (t.teardown || t.ac) && (t.teardown?.(), t.ac?.abort(ae), t.teardown = d, t.ac = null, tr(t, 0), Cn(t));
}
function Mt(e) {
	if (e.effects !== null) for (let t of e.effects) t.teardown && nr(t);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/sources.js
var Nt = /* @__PURE__ */ new Set(), Pt = /* @__PURE__ */ new Map(), Ft = !1;
function It(e, t) {
	return {
		f: 0,
		v: e,
		reactions: null,
		equals: Oe,
		rv: 0,
		wv: 0
	};
}
/* @__NO_SIDE_EFFECTS__ */
function L(e, t) {
	let n = It(e, t);
	return Vn(n), n;
}
/* @__NO_SIDE_EFFECTS__ */
function Lt(e, t = !1, n = !0) {
	let r = It(e);
	return t || (r.equals = Ae), Me && n && N !== null && N.l !== null && (N.l.s ??= []).push(r), r;
}
function R(e, t, r = !1) {
	return U !== null && (!Ln || U.f & 131072) && Ie() && U.f & 4325394 && (Bn === null || !n.call(Bn, e)) && ge(), Rt(e, r ? z(t) : t, rt);
}
function Rt(e, t, n = null) {
	if (!e.equals(t)) {
		var r = e.v;
		Fn ? Pt.set(e, t) : Pt.set(e, r), e.v = t;
		var i = ot.ensure();
		if (i.capture(e, r), e.f & 2) {
			let t = e;
			e.f & 2048 && kt(t), We(t);
		}
		e.wv = Xn(), Vt(e, g, n), Ie() && W !== null && W.f & 1024 && !(W.f & 96) && (Wn === null ? Gn([e]) : Wn.push(e)), !i.is_fork && Nt.size > 0 && !Ft && zt();
	}
	return t;
}
function zt() {
	Ft = !1;
	for (let e of Nt) e.f & 1024 && P(e, _), Zn(e) && nr(e);
	Nt.clear();
}
function Bt(e) {
	R(e, e.v + 1);
}
function Vt(e, t, n) {
	var r = e.reactions;
	if (r !== null) for (var i = Ie(), a = r.length, o = 0; o < a; o++) {
		var s = r[o], c = s.f;
		if (!(!i && s === W)) {
			var l = (c & g) === 0;
			if (l && P(s, t), c & 2) {
				var u = s;
				Qe?.delete(u), c & 65536 || (c & 512 && (s.f |= E), Vt(u, _, n));
			} else if (l) {
				var d = s;
				c & 16 && lt !== null && lt.add(d), n === null ? pt(d) : n.push(d);
			}
		}
	}
}
function z(t) {
	if (typeof t != "object" || !t || ne in t) return t;
	let n = l(t);
	if (n !== s && n !== c) return t;
	var r = /* @__PURE__ */ new Map(), i = e(t), o = /* @__PURE__ */ L(0), u = null, d = Jn, f = (e) => {
		if (Jn === d) return e();
		var t = U, n = Jn;
		Rn(null), Yn(d);
		var r = e();
		return Rn(t), Yn(n), r;
	};
	return i && r.set("length", /* @__PURE__ */ L(t.length, u)), new Proxy(t, {
		defineProperty(e, t, n) {
			(!("value" in n) || n.configurable === !1 || n.enumerable === !1 || n.writable === !1) && me();
			var i = r.get(t);
			return i === void 0 ? f(() => {
				var e = /* @__PURE__ */ L(n.value, u);
				return r.set(t, e), e;
			}) : R(i, n.value, !0), !0;
		},
		deleteProperty(e, t) {
			var n = r.get(t);
			if (n === void 0) {
				if (t in e) {
					let e = f(() => /* @__PURE__ */ L(O, u));
					r.set(t, e), Bt(o);
				}
			} else R(n, O), Bt(o);
			return !0;
		},
		get(e, n, i) {
			if (n === ne) return t;
			var o = r.get(n), s = n in e;
			if (o === void 0 && (!s || a(e, n)?.writable) && (o = f(() => /* @__PURE__ */ L(z(s ? e[n] : O), u)), r.set(n, o)), o !== void 0) {
				var c = G(o);
				return c === O ? void 0 : c;
			}
			return Reflect.get(e, n, i);
		},
		getOwnPropertyDescriptor(e, t) {
			var n = Reflect.getOwnPropertyDescriptor(e, t);
			if (n && "value" in n) {
				var i = r.get(t);
				i && (n.value = G(i));
			} else if (n === void 0) {
				var a = r.get(t), o = a?.v;
				if (a !== void 0 && o !== O) return {
					enumerable: !0,
					configurable: !0,
					value: o,
					writable: !0
				};
			}
			return n;
		},
		has(e, t) {
			if (t === ne) return !0;
			var n = r.get(t), i = n !== void 0 && n.v !== O || Reflect.has(e, t);
			return (n !== void 0 || W !== null && (!i || a(e, t)?.writable)) && (n === void 0 && (n = f(() => /* @__PURE__ */ L(i ? z(e[t]) : O, u)), r.set(t, n)), G(n) === O) ? !1 : i;
		},
		set(e, t, n, s) {
			var c = r.get(t), l = t in e;
			if (i && t === "length") for (var d = n; d < c.v; d += 1) {
				var p = r.get(d + "");
				p === void 0 ? d in e && (p = f(() => /* @__PURE__ */ L(O, u)), r.set(d + "", p)) : R(p, O);
			}
			if (c === void 0) (!l || a(e, t)?.writable) && (c = f(() => /* @__PURE__ */ L(void 0, u)), R(c, z(n)), r.set(t, c));
			else {
				l = c.v !== O;
				var m = f(() => z(n));
				R(c, m);
			}
			var h = Reflect.getOwnPropertyDescriptor(e, t);
			if (h?.set && h.set.call(s, n), !l) {
				if (i && typeof t == "string") {
					var g = r.get("length"), _ = Number(t);
					Number.isInteger(_) && _ >= g.v && R(g, _ + 1);
				}
				Bt(o);
			}
			return !0;
		},
		ownKeys(e) {
			G(o);
			var t = Reflect.ownKeys(e).filter((e) => {
				var t = r.get(e);
				return t === void 0 || t.v !== O;
			});
			for (var [n, i] of r) i.v !== O && !(n in e) && t.push(n);
			return t;
		},
		setPrototypeOf() {
			he();
		}
	});
}
function Ht(e) {
	try {
		if (typeof e == "object" && e && ne in e) return e[ne];
	} catch {}
	return e;
}
function Ut(e, t) {
	return Object.is(Ht(e), Ht(t));
}
var Wt, Gt, Kt, qt;
function Jt() {
	if (Wt === void 0) {
		Wt = window, document, Gt = /Firefox/.test(navigator.userAgent);
		var e = Element.prototype, t = Node.prototype, n = Text.prototype;
		Kt = a(t, "firstChild").get, qt = a(t, "nextSibling").get, u(e) && (e.__click = void 0, e.__className = void 0, e.__attributes = null, e.__style = void 0, e.__e = void 0), u(n) && (n.__t = void 0);
	}
}
function Yt(e = "") {
	return document.createTextNode(e);
}
/* @__NO_SIDE_EFFECTS__ */
function Xt(e) {
	return Kt.call(e);
}
/* @__NO_SIDE_EFFECTS__ */
function Zt(e) {
	return qt.call(e);
}
function B(e, t) {
	if (!k) return /* @__PURE__ */ Xt(e);
	var n = /* @__PURE__ */ Xt(A);
	if (n === null) n = A.appendChild(Yt());
	else if (t && n.nodeType !== 3) {
		var r = Yt();
		return n?.before(r), j(r), r;
	}
	return t && nn(n), j(n), n;
}
function Qt(e, t = !1) {
	if (!k) {
		var n = /* @__PURE__ */ Xt(e);
		return n instanceof Comment && n.data === "" ? /* @__PURE__ */ Zt(n) : n;
	}
	if (t) {
		if (A?.nodeType !== 3) {
			var r = Yt();
			return A?.before(r), j(r), r;
		}
		nn(A);
	}
	return A;
}
function V(e, t = 1, n = !1) {
	let r = k ? A : e;
	for (var i; t--;) i = r, r = /* @__PURE__ */ Zt(r);
	if (!k) return r;
	if (n) {
		if (r?.nodeType !== 3) {
			var a = Yt();
			return r === null ? i?.after(a) : r.before(a), j(a), a;
		}
		nn(r);
	}
	return j(r), r;
}
function $t(e) {
	e.textContent = "";
}
function en() {
	return !je || lt !== null ? !1 : (W.f & b) !== 0;
}
function tn(e, t, n) {
	let r = n ? { is: n } : void 0;
	return document.createElementNS(t ?? "http://www.w3.org/1999/xhtml", e, r);
}
function nn(e) {
	if (e.nodeValue.length < 65536) return;
	let t = e.nextSibling;
	for (; t !== null && t.nodeType === 3;) t.remove(), e.nodeValue += t.nodeValue, t = e.nextSibling;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/misc.js
function rn(e, t) {
	if (t) {
		let t = document.body;
		e.autofocus = !0, ze(() => {
			document.activeElement === t && e.focus();
		});
	}
}
var an = !1;
function on() {
	an || (an = !0, document.addEventListener("reset", (e) => {
		Promise.resolve().then(() => {
			if (!e.defaultPrevented) for (let t of e.target.elements) t.__on_r?.();
		});
	}, { capture: !0 }));
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/shared.js
function sn(e) {
	var t = U, n = W;
	Rn(null), zn(null);
	try {
		return e();
	} finally {
		Rn(t), zn(n);
	}
}
function cn(e, t, n, r = n) {
	e.addEventListener(t, () => sn(n));
	let i = e.__on_r;
	i ? e.__on_r = () => {
		i(), r(!0);
	} : e.__on_r = () => r(!0), on();
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/effects.js
function ln(e) {
	W === null && (U === null && de(e), ue()), Fn && le(e);
}
function un(e, t) {
	var n = t.last;
	n === null ? t.last = t.first = e : (n.next = e, e.prev = n, t.last = e);
}
function dn(e, t) {
	var n = W;
	n !== null && n.f & 8192 && (e |= v);
	var r = {
		ctx: N,
		deps: null,
		nodes: null,
		f: e | g | 512,
		first: null,
		fn: t,
		last: null,
		next: null,
		parent: n,
		b: n && n.b,
		prev: null,
		teardown: null,
		wv: 0,
		ac: null
	}, i = r;
	if (e & 4) nt === null ? ot.ensure().schedule(r) : nt.push(r);
	else if (t !== null) {
		try {
			nr(r);
		} catch (e) {
			throw Tn(r), e;
		}
		i.deps === null && i.teardown === null && i.nodes === null && i.first === i.last && !(i.f & 524288) && (i = i.first, e & 16 && e & 65536 && i !== null && (i.f |= S));
	}
	if (i !== null && (i.parent = n, n !== null && un(i, n), U !== null && U.f & 2 && !(e & 64))) {
		var a = U;
		(a.effects ??= []).push(i);
	}
	return r;
}
function fn() {
	return U !== null && !Ln;
}
function pn(e) {
	let t = dn(8, null);
	return P(t, h), t.teardown = e, t;
}
function mn(e) {
	ln("$effect");
	var t = W.f;
	if (!U && t & 32 && !(t & 32768)) {
		var n = N;
		(n.e ??= []).push(e);
	} else return hn(e);
}
function hn(e) {
	return dn(4 | w, e);
}
function gn(e) {
	ot.ensure();
	let t = dn(64 | C, e);
	return (e = {}) => new Promise((n) => {
		e.outro ? On(t, () => {
			Tn(t), n(void 0);
		}) : (Tn(t), n(void 0));
	});
}
function _n(e) {
	return dn(4, e);
}
function vn(e) {
	return dn(ee | C, e);
}
function yn(e, t = 0) {
	return dn(8 | t, e);
}
function H(e, t = [], n = [], r = []) {
	bt(r, t, n, (t) => {
		dn(8, () => e(...t.map(G)));
	});
}
function bn(e, t = 0) {
	return dn(16 | t, e);
}
function xn(e) {
	return dn(32 | C, e);
}
function Sn(e) {
	var t = e.teardown;
	if (t !== null) {
		let e = Fn, n = U;
		In(!0), Rn(null);
		try {
			t.call(null);
		} finally {
			In(e), Rn(n);
		}
	}
}
function Cn(e, t = !1) {
	var n = e.first;
	for (e.first = e.last = null; n !== null;) {
		let e = n.ac;
		e !== null && sn(() => {
			e.abort(ae);
		});
		var r = n.next;
		n.f & 64 ? n.parent = null : Tn(n, t), n = r;
	}
}
function wn(e) {
	for (var t = e.first; t !== null;) {
		var n = t.next;
		t.f & 32 || Tn(t), t = n;
	}
}
function Tn(e, t = !0) {
	var n = !1;
	(t || e.f & 262144) && e.nodes !== null && e.nodes.end !== null && (En(e.nodes.start, e.nodes.end), n = !0), P(e, x), Cn(e, t && !n), tr(e, 0);
	var r = e.nodes && e.nodes.t;
	if (r !== null) for (let e of r) e.stop();
	Sn(e), e.f ^= x, e.f |= y;
	var i = e.parent;
	i !== null && i.first !== null && Dn(e), e.next = e.prev = e.teardown = e.ctx = e.deps = e.fn = e.nodes = e.ac = null;
}
function En(e, t) {
	for (; e !== null;) {
		var n = e === t ? null : /* @__PURE__ */ Zt(e);
		e.remove(), e = n;
	}
}
function Dn(e) {
	var t = e.parent, n = e.prev, r = e.next;
	n !== null && (n.next = r), r !== null && (r.prev = n), t !== null && (t.first === e && (t.first = r), t.last === e && (t.last = n));
}
function On(e, t, n = !0) {
	var r = [];
	kn(e, r, !0);
	var i = () => {
		n && Tn(e), t && t();
	}, a = r.length;
	if (a > 0) {
		var o = () => --a || i();
		for (var s of r) s.out(o);
	} else i();
}
function kn(e, t, n) {
	if (!(e.f & 8192)) {
		e.f ^= v;
		var r = e.nodes && e.nodes.t;
		if (r !== null) for (let e of r) (e.is_global || n) && t.push(e);
		for (var i = e.first; i !== null;) {
			var a = i.next, o = (i.f & 65536) != 0 || (i.f & 32) != 0 && (e.f & 16) != 0;
			kn(i, t, o ? n : !1), i = a;
		}
	}
}
function An(e) {
	jn(e, !0);
}
function jn(e, t) {
	if (e.f & 8192) {
		e.f ^= v, e.f & 1024 || (P(e, g), ot.ensure().schedule(e));
		for (var n = e.first; n !== null;) {
			var r = n.next, i = (n.f & 65536) != 0 || (n.f & 32) != 0;
			jn(n, i ? t : !1), n = r;
		}
		var a = e.nodes && e.nodes.t;
		if (a !== null) for (let e of a) (e.is_global || t) && e.in();
	}
}
function Mn(e, t) {
	if (e.nodes) for (var n = e.nodes.start, r = e.nodes.end; n !== null;) {
		var i = n === r ? null : /* @__PURE__ */ Zt(n);
		t.append(n), n = i;
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/legacy.js
var Nn = null, Pn = !1, Fn = !1;
function In(e) {
	Fn = e;
}
var U = null, Ln = !1;
function Rn(e) {
	U = e;
}
var W = null;
function zn(e) {
	W = e;
}
var Bn = null;
function Vn(e) {
	U !== null && (!je || U.f & 2) && (Bn === null ? Bn = [e] : Bn.push(e));
}
var Hn = null, Un = 0, Wn = null;
function Gn(e) {
	Wn = e;
}
var Kn = 1, qn = 0, Jn = qn;
function Yn(e) {
	Jn = e;
}
function Xn() {
	return ++Kn;
}
function Zn(e) {
	var t = e.f;
	if (t & 2048) return !0;
	if (t & 2 && (e.f &= ~E), t & 4096) {
		for (var n = e.deps, r = n.length, i = 0; i < r; i++) {
			var a = n[i];
			if (Zn(a) && At(a), a.wv > e.wv) return !0;
		}
		t & 512 && Qe === null && P(e, h);
	}
	return !1;
}
function Qn(e, t, r = !0) {
	var i = e.reactions;
	if (i !== null && !(!je && Bn !== null && n.call(Bn, e))) for (var a = 0; a < i.length; a++) {
		var o = i[a];
		o.f & 2 ? Qn(o, t, !1) : t === o && (r ? P(o, g) : o.f & 1024 && P(o, _), pt(o));
	}
}
function $n(e) {
	var t = Hn, n = Un, r = Wn, i = U, a = Bn, o = N, s = Ln, c = Jn, l = e.f;
	Hn = null, Un = 0, Wn = null, U = l & 96 ? null : e, Bn = null, Ne(e.ctx), Ln = !1, Jn = ++qn, e.ac !== null && (sn(() => {
		e.ac.abort(ae);
	}), e.ac = null);
	try {
		e.f |= D;
		var u = e.fn, d = u();
		e.f |= b;
		var f = e.deps, p = F?.is_fork;
		if (Hn !== null) {
			var m;
			if (p || tr(e, Un), f !== null && Un > 0) for (f.length = Un + Hn.length, m = 0; m < Hn.length; m++) f[Un + m] = Hn[m];
			else e.deps = f = Hn;
			if (fn() && e.f & 512) for (m = Un; m < f.length; m++) (f[m].reactions ??= []).push(e);
		} else !p && f !== null && Un < f.length && (tr(e, Un), f.length = Un);
		if (Ie() && Wn !== null && !Ln && f !== null && !(e.f & 6146)) for (m = 0; m < Wn.length; m++) Qn(Wn[m], e);
		if (i !== null && i !== e) {
			if (qn++, i.deps !== null) for (let e = 0; e < n; e += 1) i.deps[e].rv = qn;
			if (t !== null) for (let e of t) e.rv = qn;
			Wn !== null && (r === null ? r = Wn : r.push(...Wn));
		}
		return e.f & 8388608 && (e.f ^= te), d;
	} catch (e) {
		return Ve(e);
	} finally {
		e.f ^= D, Hn = t, Un = n, Wn = r, U = i, Bn = a, Ne(o), Ln = s, Jn = c;
	}
}
function er(e, r) {
	let i = r.reactions;
	if (i !== null) {
		var a = t.call(i, e);
		if (a !== -1) {
			var o = i.length - 1;
			o === 0 ? i = r.reactions = null : (i[a] = i[o], i.pop());
		}
	}
	if (i === null && r.f & 2 && (Hn === null || !n.call(Hn, r))) {
		var s = r;
		s.f & 512 && (s.f ^= 512, s.f &= ~E), We(s), jt(s), tr(s, 0);
	}
}
function tr(e, t) {
	var n = e.deps;
	if (n !== null) for (var r = t; r < n.length; r++) er(e, n[r]);
}
function nr(e) {
	var t = e.f;
	if (!(t & 16384)) {
		P(e, h);
		var n = W, r = Pn;
		W = e, Pn = !0;
		try {
			t & 16777232 ? wn(e) : Cn(e), Sn(e);
			var i = $n(e);
			e.teardown = typeof i == "function" ? i : null, e.wv = Kn;
		} finally {
			Pn = r, W = n;
		}
	}
}
async function rr() {
	if (je) return new Promise((e) => {
		requestAnimationFrame(() => e()), setTimeout(() => e());
	});
	await Promise.resolve(), st();
}
function G(e) {
	var t = (e.f & 2) != 0;
	if (Nn?.add(e), U !== null && !Ln && !(W !== null && W.f & 16384) && (Bn === null || !n.call(Bn, e))) {
		var r = U.deps;
		if (U.f & 2097152) e.rv < qn && (e.rv = qn, Hn === null && r !== null && r[Un] === e ? Un++ : Hn === null ? Hn = [e] : Hn.push(e));
		else {
			(U.deps ??= []).push(e);
			var i = e.reactions;
			i === null ? e.reactions = [U] : n.call(i, U) || i.push(U);
		}
	}
	if (Fn && Pt.has(e)) return Pt.get(e);
	if (t) {
		var a = e;
		if (Fn) {
			var o = a.v;
			return (!(a.f & 1024) && a.reactions !== null || ar(a)) && (o = kt(a)), Pt.set(a, o), o;
		}
		var s = (a.f & 512) == 0 && !Ln && U !== null && (Pn || (U.f & 512) != 0), c = (a.f & b) === 0;
		Zn(a) && (s && (a.f |= 512), At(a)), s && !c && (Mt(a), ir(a));
	}
	if (Qe?.has(e)) return Qe.get(e);
	if (e.f & 8388608) throw e.v;
	return e.v;
}
function ir(e) {
	if (e.f |= 512, e.deps !== null) for (let t of e.deps) (t.reactions ??= []).push(e), t.f & 2 && !(t.f & 512) && (Mt(t), ir(t));
}
function ar(e) {
	if (e.v === O) return !0;
	if (e.deps === null) return !1;
	for (let t of e.deps) if (Pt.has(t) || t.f & 2 && ar(t)) return !0;
	return !1;
}
function or(e) {
	var t = Ln;
	try {
		return Ln = !0, e();
	} finally {
		Ln = t;
	}
}
[.../* @__PURE__ */ "allowfullscreen.async.autofocus.autoplay.checked.controls.default.disabled.formnovalidate.indeterminate.inert.ismap.loop.multiple.muted.nomodule.novalidate.open.playsinline.readonly.required.reversed.seamless.selected.webkitdirectory.defer.disablepictureinpicture.disableremoteplayback".split(".")];
var sr = ["touchstart", "touchmove"];
function cr(e) {
	return sr.includes(e);
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/events.js
var lr = Symbol("events"), ur = /* @__PURE__ */ new Set(), dr = /* @__PURE__ */ new Set();
function fr(e, t, n, r = {}) {
	function i(e) {
		if (r.capture || gr.call(t, e), !e.cancelBubble) return sn(() => n?.call(this, e));
	}
	return e.startsWith("pointer") || e.startsWith("touch") || e === "wheel" ? ze(() => {
		t.addEventListener(e, i, r);
	}) : t.addEventListener(e, i, r), i;
}
function pr(e, t, n, r, i) {
	var a = {
		capture: r,
		passive: i
	}, o = fr(e, t, n, a);
	(t === document.body || t === window || t === document || t instanceof HTMLMediaElement) && pn(() => {
		t.removeEventListener(e, o, a);
	});
}
function K(e, t, n) {
	(t[lr] ??= {})[e] = n;
}
function mr(e) {
	for (var t = 0; t < e.length; t++) ur.add(e[t]);
	for (var n of dr) n(e);
}
var hr = null;
function gr(e) {
	var t = this, n = t.ownerDocument, r = e.type, a = e.composedPath?.() || [], o = a[0] || e.target;
	hr = e;
	var s = 0, c = hr === e && e[lr];
	if (c) {
		var l = a.indexOf(c);
		if (l !== -1 && (t === document || t === window)) {
			e[lr] = t;
			return;
		}
		var u = a.indexOf(t);
		if (u === -1) return;
		l <= u && (s = l);
	}
	if (o = a[s] || e.target, o !== t) {
		i(e, "currentTarget", {
			configurable: !0,
			get() {
				return o || n;
			}
		});
		var d = U, f = W;
		Rn(null), zn(null);
		try {
			for (var p, m = []; o !== null;) {
				var h = o.assignedSlot || o.parentNode || o.host || null;
				try {
					var g = o[lr]?.[r];
					g != null && (!o.disabled || e.target === o) && g.call(o, e);
				} catch (e) {
					p ? m.push(e) : p = e;
				}
				if (e.cancelBubble || h === t || h === null) break;
				o = h;
			}
			if (p) {
				for (let e of m) queueMicrotask(() => {
					throw e;
				});
				throw p;
			}
		} finally {
			e[lr] = t, delete e.currentTarget, Rn(d), zn(f);
		}
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/reconciler.js
var _r = globalThis?.window?.trustedTypes && /* @__PURE__ */ globalThis.window.trustedTypes.createPolicy("svelte-trusted-html", { createHTML: (e) => e });
function vr(e) {
	return _r?.createHTML(e) ?? e;
}
function yr(e) {
	var t = tn("template");
	return t.innerHTML = vr(e.replaceAll("<!>", "<!---->")), t.content;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/template.js
function br(e, t) {
	var n = W;
	n.nodes === null && (n.nodes = {
		start: e,
		end: t,
		a: null,
		t: null
	});
}
/* @__NO_SIDE_EFFECTS__ */
function q(e, t) {
	var n = (t & 1) != 0, r = (t & 2) != 0, i, a = !e.startsWith("<!>");
	return () => {
		if (k) return br(A, null), A;
		i === void 0 && (i = yr(a ? e : "<!>" + e), n || (i = /* @__PURE__ */ Xt(i)));
		var t = r || Gt ? document.importNode(i, !0) : i.cloneNode(!0);
		if (n) {
			var o = /* @__PURE__ */ Xt(t), s = t.lastChild;
			br(o, s);
		} else br(t, t);
		return t;
	};
}
function xr() {
	if (k) return br(A, null), A;
	var e = document.createDocumentFragment(), t = document.createComment(""), n = Yt();
	return e.append(t, n), br(t, n), e;
}
function J(e, t) {
	if (k) {
		var n = W;
		(!(n.f & 32768) || n.nodes.end === null) && (n.nodes.end = A), we();
		return;
	}
	e !== null && e.before(t);
}
function Y(e, t) {
	var n = t == null ? "" : typeof t == "object" ? `${t}` : t;
	n !== (e.__t ??= e.nodeValue) && (e.__t = n, e.nodeValue = `${n}`);
}
function Sr(e, t) {
	return wr(e, t);
}
var Cr = /* @__PURE__ */ new Map();
function wr(e, { target: t, anchor: n, props: i = {}, events: a, context: o, intro: s = !0, transformError: c }) {
	Jt();
	var l = void 0, u = gn(() => {
		var s = n ?? t.appendChild(Yt());
		vt(s, { pending: () => {} }, (t) => {
			Pe({});
			var n = N;
			if (o && (n.c = o), a && (i.$$events = a), k && br(t, null), l = e(t, i) || {}, k && (W.nodes.end = A, A === null || A.nodeType !== 8 || A.data !== "]")) throw be(), ve;
			Fe();
		}, c);
		var u = /* @__PURE__ */ new Set(), d = (e) => {
			for (var n = 0; n < e.length; n++) {
				var r = e[n];
				if (!u.has(r)) {
					u.add(r);
					var i = cr(r);
					for (let e of [t, document]) {
						var a = Cr.get(e);
						a === void 0 && (a = /* @__PURE__ */ new Map(), Cr.set(e, a));
						var o = a.get(r);
						o === void 0 ? (e.addEventListener(r, gr, { passive: i }), a.set(r, 1)) : a.set(r, o + 1);
					}
				}
			}
		};
		return d(r(ur)), dr.add(d), () => {
			for (var e of u) for (let n of [t, document]) {
				var r = Cr.get(n), i = r.get(e);
				--i == 0 ? (n.removeEventListener(e, gr), r.delete(e), r.size === 0 && Cr.delete(n)) : r.set(e, i);
			}
			dr.delete(d), s !== n && s.parentNode?.removeChild(s);
		};
	});
	return Tr.set(l, u), l;
}
var Tr = /* @__PURE__ */ new WeakMap(), Er = class {
	anchor;
	#e = /* @__PURE__ */ new Map();
	#t = /* @__PURE__ */ new Map();
	#n = /* @__PURE__ */ new Map();
	#r = /* @__PURE__ */ new Set();
	#i = !0;
	constructor(e, t = !0) {
		this.anchor = e, this.#i = t;
	}
	#a = (e) => {
		if (this.#e.has(e)) {
			var t = this.#e.get(e), n = this.#t.get(t);
			if (n) An(n), this.#r.delete(t);
			else {
				var r = this.#n.get(t);
				r && (this.#t.set(t, r.effect), this.#n.delete(t), r.fragment.lastChild.remove(), this.anchor.before(r.fragment), n = r.effect);
			}
			for (let [t, n] of this.#e) {
				if (this.#e.delete(t), t === e) break;
				let r = this.#n.get(n);
				r && (Tn(r.effect), this.#n.delete(n));
			}
			for (let [e, r] of this.#t) {
				if (e === t || this.#r.has(e)) continue;
				let i = () => {
					if (Array.from(this.#e.values()).includes(e)) {
						var t = document.createDocumentFragment();
						Mn(r, t), t.append(Yt()), this.#n.set(e, {
							effect: r,
							fragment: t
						});
					} else Tn(r);
					this.#r.delete(e), this.#t.delete(e);
				};
				this.#i || !n ? (this.#r.add(e), On(r, i, !1)) : i();
			}
		}
	};
	#o = (e) => {
		this.#e.delete(e);
		let t = Array.from(this.#e.values());
		for (let [e, n] of this.#n) t.includes(e) || (Tn(n.effect), this.#n.delete(e));
	};
	ensure(e, t) {
		var n = F, r = en();
		if (t && !this.#t.has(e) && !this.#n.has(e)) if (r) {
			var i = document.createDocumentFragment(), a = Yt();
			i.append(a), this.#n.set(e, {
				effect: xn(() => t(a)),
				fragment: i
			});
		} else this.#t.set(e, xn(() => t(this.anchor)));
		if (this.#e.set(n, e), r) {
			for (let [t, r] of this.#t) t === e ? n.unskip_effect(r) : n.skip_effect(r);
			for (let [t, r] of this.#n) t === e ? n.unskip_effect(r.effect) : n.skip_effect(r.effect);
			n.oncommit(this.#a), n.ondiscard(this.#o);
		} else k && (this.anchor = A), this.#a(n);
	}
};
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/if.js
function X(e, t, n = !1) {
	var r;
	k && (r = A, we());
	var i = new Er(e), a = n ? S : 0;
	function o(e, t) {
		if (k) {
			var n = De(r);
			if (e !== parseInt(n.substring(1))) {
				var a = Ee();
				j(a), i.anchor = a, Ce(!1), i.ensure(e, t), Ce(!0);
				return;
			}
		}
		i.ensure(e, t);
	}
	bn(() => {
		var e = !1;
		t((t, n = 0) => {
			e = !0, o(n, t);
		}), e || o(-1, null);
	}, a);
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/each.js
function Dr(e, t) {
	return t;
}
function Or(e, t, n) {
	for (var i = [], a = t.length, o, s = t.length, c = 0; c < a; c++) {
		let n = t[c];
		On(n, () => {
			if (o) {
				if (o.pending.delete(n), o.done.add(n), o.pending.size === 0) {
					var t = e.outrogroups;
					kr(e, r(o.done)), t.delete(o), t.size === 0 && (e.outrogroups = null);
				}
			} else --s;
		}, !1);
	}
	if (s === 0) {
		var l = i.length === 0 && n !== null;
		if (l) {
			var u = n, d = u.parentNode;
			$t(d), d.append(u), e.items.clear();
		}
		kr(e, t, !l);
	} else o = {
		pending: new Set(t),
		done: /* @__PURE__ */ new Set()
	}, (e.outrogroups ??= /* @__PURE__ */ new Set()).add(o);
}
function kr(e, t, n = !0) {
	var r;
	if (e.pending.size > 0) {
		r = /* @__PURE__ */ new Set();
		for (let t of e.pending.values()) for (let n of t) r.add(e.items.get(n).e);
	}
	for (var i = 0; i < t.length; i++) {
		var a = t[i];
		r?.has(a) ? (a.f |= T, Mn(a, document.createDocumentFragment())) : Tn(t[i], n);
	}
}
var Ar;
function Z(t, n, i, a, o, s = null) {
	var c = t, l = /* @__PURE__ */ new Map();
	if (n & 4) {
		var u = t;
		c = k ? j(/* @__PURE__ */ Xt(u)) : u.appendChild(Yt());
	}
	k && we();
	var d = null, f = /* @__PURE__ */ Et(() => {
		var t = i();
		return e(t) ? t : t == null ? [] : r(t);
	}), p, m = /* @__PURE__ */ new Map(), h = !0;
	function g(e) {
		v.effect.f & 16384 || (v.pending.delete(e), v.fallback = d, Mr(v, p, c, n, a), d !== null && (p.length === 0 ? d.f & 33554432 ? (d.f ^= T, Pr(d, null, c)) : An(d) : On(d, () => {
			d = null;
		})));
	}
	function _(e) {
		v.pending.delete(e);
	}
	var v = {
		effect: bn(() => {
			p = G(f);
			var e = p.length;
			let t = !1;
			k && De(c) === "[!" != (e === 0) && (c = Ee(), j(c), Ce(!1), t = !0);
			for (var r = /* @__PURE__ */ new Set(), u = F, v = en(), y = 0; y < e; y += 1) {
				k && A.nodeType === 8 && A.data === "]" && (c = A, t = !0, Ce(!1));
				var b = p[y], x = a(b, y), S = h ? null : l.get(x);
				S ? (S.v && Rt(S.v, b), S.i && Rt(S.i, y), v && u.unskip_effect(S.e)) : (S = Nr(l, h ? c : Ar ??= Yt(), b, x, y, o, n, i), h || (S.e.f |= T), l.set(x, S)), r.add(x);
			}
			if (e === 0 && s && !d && (h ? d = xn(() => s(c)) : (d = xn(() => s(Ar ??= Yt())), d.f |= T)), e > r.size && ce("", "", ""), k && e > 0 && j(Ee()), !h) if (m.set(u, r), v) {
				for (let [e, t] of l) r.has(e) || u.skip_effect(t.e);
				u.oncommit(g), u.ondiscard(_);
			} else g(u);
			t && Ce(!0), G(f);
		}),
		flags: n,
		items: l,
		pending: m,
		outrogroups: null,
		fallback: d
	};
	h = !1, k && (c = A);
}
function jr(e) {
	for (; e !== null && !(e.f & 32);) e = e.next;
	return e;
}
function Mr(e, t, n, i, a) {
	var o = (i & 8) != 0, s = t.length, c = e.items, l = jr(e.effect.first), u, d = null, f, p = [], m = [], h, g, _, v;
	if (o) for (v = 0; v < s; v += 1) h = t[v], g = a(h, v), _ = c.get(g).e, _.f & 33554432 || (_.nodes?.a?.measure(), (f ??= /* @__PURE__ */ new Set()).add(_));
	for (v = 0; v < s; v += 1) {
		if (h = t[v], g = a(h, v), _ = c.get(g).e, e.outrogroups !== null) for (let t of e.outrogroups) t.pending.delete(_), t.done.delete(_);
		if (_.f & 33554432) if (_.f ^= T, _ === l) Pr(_, null, n);
		else {
			var y = d ? d.next : l;
			_ === e.effect.last && (e.effect.last = _.prev), _.prev && (_.prev.next = _.next), _.next && (_.next.prev = _.prev), Fr(e, d, _), Fr(e, _, y), Pr(_, y, n), d = _, p = [], m = [], l = jr(d.next);
			continue;
		}
		if (_.f & 8192 && (An(_), o && (_.nodes?.a?.unfix(), (f ??= /* @__PURE__ */ new Set()).delete(_))), _ !== l) {
			if (u !== void 0 && u.has(_)) {
				if (p.length < m.length) {
					var b = m[0], x;
					d = b.prev;
					var S = p[0], C = p[p.length - 1];
					for (x = 0; x < p.length; x += 1) Pr(p[x], b, n);
					for (x = 0; x < m.length; x += 1) u.delete(m[x]);
					Fr(e, S.prev, C.next), Fr(e, d, S), Fr(e, C, b), l = b, d = C, --v, p = [], m = [];
				} else u.delete(_), Pr(_, l, n), Fr(e, _.prev, _.next), Fr(e, _, d === null ? e.effect.first : d.next), Fr(e, d, _), d = _;
				continue;
			}
			for (p = [], m = []; l !== null && l !== _;) (u ??= /* @__PURE__ */ new Set()).add(l), m.push(l), l = jr(l.next);
			if (l === null) continue;
		}
		_.f & 33554432 || p.push(_), d = _, l = jr(_.next);
	}
	if (e.outrogroups !== null) {
		for (let t of e.outrogroups) t.pending.size === 0 && (kr(e, r(t.done)), e.outrogroups?.delete(t));
		e.outrogroups.size === 0 && (e.outrogroups = null);
	}
	if (l !== null || u !== void 0) {
		var w = [];
		if (u !== void 0) for (_ of u) _.f & 8192 || w.push(_);
		for (; l !== null;) !(l.f & 8192) && l !== e.fallback && w.push(l), l = jr(l.next);
		var E = w.length;
		if (E > 0) {
			var D = i & 4 && s === 0 ? n : null;
			if (o) {
				for (v = 0; v < E; v += 1) w[v].nodes?.a?.measure();
				for (v = 0; v < E; v += 1) w[v].nodes?.a?.fix();
			}
			Or(e, w, D);
		}
	}
	o && ze(() => {
		if (f !== void 0) for (_ of f) _.nodes?.a?.apply();
	});
}
function Nr(e, t, n, r, i, a, o, s) {
	var c = o & 1 ? o & 16 ? It(n) : /* @__PURE__ */ Lt(n, !1, !1) : null, l = o & 2 ? It(i) : null;
	return {
		v: c,
		i: l,
		e: xn(() => (a(t, c ?? n, l ?? i, s), () => {
			e.delete(r);
		}))
	};
}
function Pr(e, t, n) {
	if (e.nodes) for (var r = e.nodes.start, i = e.nodes.end, a = t && !(t.f & 33554432) ? t.nodes.start : n; r !== null;) {
		var o = /* @__PURE__ */ Zt(r);
		if (a.before(r), r === i) return;
		r = o;
	}
}
function Fr(e, t, n) {
	t === null ? e.effect.first = n : t.next = n, n === null ? e.effect.last = t : n.prev = t;
}
//#endregion
//#region node_modules/svelte/src/internal/shared/attributes.js
var Ir = [..." 	\n\r\f\xA0\v﻿"];
function Lr(e, t, n) {
	var r = e == null ? "" : "" + e;
	if (t && (r = r ? r + " " + t : t), n) {
		for (var i of Object.keys(n)) if (n[i]) r = r ? r + " " + i : i;
		else if (r.length) for (var a = i.length, o = 0; (o = r.indexOf(i, o)) >= 0;) {
			var s = o + a;
			(o === 0 || Ir.includes(r[o - 1])) && (s === r.length || Ir.includes(r[s])) ? r = (o === 0 ? "" : r.substring(0, o)) + r.substring(s + 1) : o = s;
		}
	}
	return r === "" ? null : r;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/class.js
function Rr(e, t, n, r, i, a) {
	var o = e.__className;
	if (k || o !== n || o === void 0) {
		var s = Lr(n, r, a);
		(!k || s !== e.getAttribute("class")) && (s == null ? e.removeAttribute("class") : t ? e.className = s : e.setAttribute("class", s)), e.__className = n;
	} else if (a && i !== a) for (var c in a) {
		var l = !!a[c];
		(i == null || l !== !!i[c]) && e.classList.toggle(c, l);
	}
	return a;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/select.js
function zr(t, n, r = !1) {
	if (t.multiple) {
		if (n == null) return;
		if (!e(n)) return xe();
		for (var i of t.options) i.selected = n.includes(Vr(i));
		return;
	}
	for (i of t.options) if (Ut(Vr(i), n)) {
		i.selected = !0;
		return;
	}
	(!r || n !== void 0) && (t.selectedIndex = -1);
}
function Br(e) {
	var t = new MutationObserver(() => {
		zr(e, e.__value);
	});
	t.observe(e, {
		childList: !0,
		subtree: !0,
		attributes: !0,
		attributeFilter: ["value"]
	}), pn(() => {
		t.disconnect();
	});
}
function Vr(e) {
	return "__value" in e ? e.__value : e.value;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/attributes.js
var Hr = Symbol("is custom element"), Ur = Symbol("is html"), Wr = oe ? "link" : "LINK", Gr = oe ? "progress" : "PROGRESS";
function Q(e) {
	if (k) {
		var t = !1, n = () => {
			if (!t) {
				if (t = !0, e.hasAttribute("value")) {
					var n = e.value;
					Jr(e, "value", null), e.value = n;
				}
				if (e.hasAttribute("checked")) {
					var r = e.checked;
					Jr(e, "checked", null), e.checked = r;
				}
			}
		};
		e.__on_r = n, ze(n), on();
	}
}
function Kr(e, t) {
	var n = Yr(e);
	n.value === (n.value = t ?? void 0) || e.value === t && (t !== 0 || e.nodeName !== Gr) || (e.value = t ?? "");
}
function qr(e, t) {
	var n = Yr(e);
	n.checked !== (n.checked = t ?? void 0) && (e.checked = t);
}
function Jr(e, t, n, r) {
	var i = Yr(e);
	k && (i[t] = e.getAttribute(t), t === "src" || t === "srcset" || t === "href" && e.nodeName === Wr) || i[t] !== (i[t] = n) && (t === "loading" && (e[ie] = n), n == null ? e.removeAttribute(t) : typeof n != "string" && Zr(e).includes(t) ? e[t] = n : e.setAttribute(t, n));
}
function Yr(e) {
	return e.__attributes ??= {
		[Hr]: e.nodeName.includes("-"),
		[Ur]: e.namespaceURI === ye
	};
}
var Xr = /* @__PURE__ */ new Map();
function Zr(e) {
	var t = e.getAttribute("is") || e.nodeName, n = Xr.get(t);
	if (n) return n;
	Xr.set(t, n = []);
	for (var r, i = e, a = Element.prototype; a !== i;) {
		for (var s in r = o(i), r) r[s].set && n.push(s);
		i = l(i);
	}
	return n;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/input.js
function Qr(e, t, n = t) {
	var r = /* @__PURE__ */ new WeakSet();
	cn(e, "input", async (i) => {
		var a = i ? e.defaultValue : e.value;
		if (a = $r(e) ? ei(a) : a, n(a), F !== null && r.add(F), await rr(), a !== (a = t())) {
			var o = e.selectionStart, s = e.selectionEnd, c = e.value.length;
			if (e.value = a ?? "", s !== null) {
				var l = e.value.length;
				o === s && s === c && l > c ? (e.selectionStart = l, e.selectionEnd = l) : (e.selectionStart = o, e.selectionEnd = Math.min(s, l));
			}
		}
	}), (k && e.defaultValue !== e.value || or(t) == null && e.value) && (n($r(e) ? ei(e.value) : e.value), F !== null && r.add(F)), yn(() => {
		var n = t();
		if (e === document.activeElement) {
			var i = je ? Ze : F;
			if (r.has(i)) return;
		}
		$r(e) && n === ei(e.value) || e.type === "date" && !n && !e.value || n !== e.value && (e.value = n ?? "");
	});
}
function $r(e) {
	var t = e.type;
	return t === "number" || t === "range";
}
function ei(e) {
	return e === "" ? null : +e;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/this.js
function ti(e, t) {
	return e === t || e?.[ne] === t;
}
function ni(e = {}, t, n, r) {
	var i = N.r, a = W;
	return _n(() => {
		var o, s;
		return yn(() => {
			o = s, s = r?.() || [], or(() => {
				e !== n(...s) && (t(e, ...s), o && ti(n(...o), e) && t(null, ...o));
			});
		}), () => {
			let r = a;
			for (; r !== i && r.parent !== null && r.parent.f & 33554432;) r = r.parent;
			let o = () => {
				s && ti(n(...s), e) && t(null, ...s);
			}, c = r.teardown;
			r.teardown = () => {
				o(), c?.();
			};
		};
	}), e;
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/props.js
function $(e, t, n, r) {
	var i = !Me || (n & 2) != 0, o = (n & 8) != 0, s = (n & 16) != 0, c = r, l = !0, u = () => (l && (l = !1, c = s ? or(r) : r), c);
	let d;
	if (o) {
		var f = ne in e || re in e;
		d = a(e, t)?.set ?? (f && t in e ? (n) => e[t] = n : void 0);
	}
	var p, m = !1;
	o ? [p, m] = Ye(() => e[t]) : p = e[t], p === void 0 && r !== void 0 && (p = u(), d && (i && pe(t), d(p)));
	var h = i ? () => {
		var n = e[t];
		return n === void 0 ? u() : (l = !0, n);
	} : () => {
		var n = e[t];
		return n !== void 0 && (c = void 0), n === void 0 ? c : n;
	};
	if (i && !(n & 4)) return h;
	if (d) {
		var g = e.$$legacy;
		return (function(e, t) {
			return arguments.length > 0 ? ((!i || !t || g || m) && d(t ? h() : e), e) : h();
		});
	}
	var _ = !1, v = (n & 1 ? wt : Et)(() => (_ = !1, h()));
	o && G(v);
	var y = W;
	return (function(e, t) {
		if (arguments.length > 0) {
			let n = t ? G(v) : i && o ? z(e) : e;
			return R(v, n), _ = !0, c !== void 0 && (c = n), e;
		}
		return Fn && _ || y.f & 16384 ? v.v : G(v);
	});
}
//#endregion
//#region node_modules/svelte/src/internal/disclose-version.js
typeof window < "u" && ((window.__svelte ??= {}).v ??= /* @__PURE__ */ new Set()).add("5");
//#endregion
//#region src/lib/markup.js
function ri(e, t, n, r) {
	let i = {
		materials: {
			global: "materials_markup",
			override: "materials",
			enabled: "materials"
		},
		labor: {
			global: "labor_markup",
			override: "labor",
			enabled: "labor"
		},
		equipment: {
			global: "equipment_markup",
			override: "equipment",
			enabled: "equipment"
		},
		subs: {
			global: "subs_markup",
			override: "subs",
			enabled: "subs"
		},
		other: {
			global: "other_markup",
			override: "other",
			enabled: "other"
		}
	}[e];
	return !i || !r[i.enabled] ? 0 : n[i.override] ?? t[i.global];
}
function ii(e, t) {
	return e * (1 + t / 100);
}
function ai(e, t, n) {
	return e * ii(t, n);
}
function oi(e, t) {
	let n = {
		materials: 0,
		labor: 0,
		equipment: 0,
		subs: 0,
		other: 0
	}, r = 0, i = 0, a = (a) => {
		let o = ri(a.category_type, t, e.markup_overrides, e.markup_enabled), s = a.quantity * a.unit_price, c = ai(a.quantity, a.unit_price, o);
		r += s, i += c, n[a.category_type] !== void 0 && (n[a.category_type] += c);
	};
	for (let t of e.line_items) a(t);
	for (let t of e.component_groups) for (let e of t.line_items) a(e);
	return i += e.lump_sum, {
		base: r,
		withMarkup: i,
		byType: n
	};
}
function si(e, t) {
	let n = {
		materials: 0,
		labor: 0,
		equipment: 0,
		subs: 0,
		other: 0
	}, r = 0, i = 0;
	for (let a of e.subcategories) {
		let e = oi(a, t);
		r += e.base, i += e.withMarkup;
		for (let t of Object.keys(n)) n[t] += e.byType[t];
	}
	return {
		base: r,
		withMarkup: i,
		byType: n
	};
}
function ci(e) {
	let t = {
		materials: 0,
		labor: 0,
		equipment: 0,
		subs: 0,
		other: 0
	}, n = 0, r = 0;
	for (let i of e.sections) {
		let a = si(i, e.globals);
		n += a.base, r += a.withMarkup;
		for (let e of Object.keys(t)) t[e] += a.byType[e];
	}
	return {
		base: n,
		withMarkup: r,
		byType: t
	};
}
function li(e) {
	return new Intl.NumberFormat("en-US", {
		style: "currency",
		currency: "USD",
		minimumFractionDigits: 2
	}).format(e);
}
function ui(e) {
	return `${e}%`;
}
//#endregion
//#region src/lib/LineItemRow.svelte
var di = /* @__PURE__ */ q("<option> </option>"), fi = /* @__PURE__ */ q("<tr class=\"border-b border-white/[0.04]\"><td colspan=\"9\" class=\"px-2 py-1\"><div class=\"flex items-center gap-2 pl-6\"><svg class=\"w-3 h-3 text-white/15 shrink-0\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M3 10h10a3 3 0 013 3v1\"></path></svg> <input type=\"text\" placeholder=\"Add a note...\" class=\"w-full px-1 py-0.5 text-xs bg-transparent border-0 border-b border-white/[0.04] text-[var(--color-concrete)]\n						focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/15 font-[var(--font-body)]\"/></div></td></tr>"), pi = /* @__PURE__ */ q("<tr class=\"border-b border-white/[0.04] hover:bg-white/[0.02] text-sm group\"><td class=\"px-1 py-1 w-24\"><select class=\"w-full text-xs px-1 py-1 border-0 bg-transparent font-medium cursor-pointer\n				focus:ring-1 focus:ring-[var(--color-sunburst)] text-[var(--color-concrete)] font-[var(--font-ui)]\" style=\"color-scheme: dark;\"></select></td><td class=\"px-1 py-1\"><input type=\"text\" class=\"w-full px-1 py-0.5 text-[var(--color-white)] bg-transparent border-0\n				focus:ring-1 focus:ring-[var(--color-sunburst)] font-[var(--font-body)]\" placeholder=\"Item name\"/></td><td class=\"px-1 py-1 w-20\"><input type=\"number\" step=\"any\" min=\"0\" class=\"w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 text-[var(--color-white)]\n				focus:ring-1 focus:ring-[var(--color-sunburst)]\"/></td><td class=\"px-1 py-1 w-16\"><input type=\"text\" class=\"w-full text-center text-white/30 px-1 py-0.5 bg-transparent border-0\n				focus:ring-1 focus:ring-[var(--color-sunburst)] text-sm font-[var(--font-body)]\" placeholder=\"ea\"/></td><td class=\"px-1 py-1 w-24\"><input type=\"text\" inputmode=\"decimal\" class=\"w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 text-[var(--color-white)]\n				focus:ring-1 focus:ring-[var(--color-sunburst)]\"/></td><td class=\"px-2 py-1.5 text-right font-mono text-white/25 w-16 text-xs\"> </td><td class=\"px-2 py-1.5 text-right font-mono text-white/40 w-24 text-xs\"> </td><td class=\"px-2 py-1.5 text-right font-mono font-medium text-[var(--color-white)] w-24\"> </td><td class=\"px-1 py-1 w-16\"><div class=\"flex items-center gap-0.5\"><button title=\"Toggle note\"><svg class=\"w-3.5 h-3.5\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z\"></path></svg></button> <button class=\"opacity-0 group-hover:opacity-100 text-white/20 hover:text-red-400 transition-opacity p-0.5\" title=\"Delete item\"><svg class=\"w-3.5 h-3.5\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></button></div></td></tr> <!>", 1);
function mi(e, t) {
	Pe(t, !0);
	let n = $(t, "item", 7), r = [
		"materials",
		"labor",
		"equipment",
		"subs",
		"other"
	], i = /* @__PURE__ */ I(() => ri(n().category_type, t.globals, t.markupOverrides, t.markupEnabled)), a = /* @__PURE__ */ I(() => ii(n().unit_price, G(i))), o = /* @__PURE__ */ I(() => ai(n().quantity, n().unit_price, G(i))), s = /* @__PURE__ */ L(!!n().description);
	function c() {
		t.onchange?.();
	}
	function l(e) {
		let t = parseFloat(e.target.value);
		isNaN(t) || (n().quantity = t, c());
	}
	function u(e) {
		let t = parseFloat(e.target.value);
		isNaN(t) || (n().unit_price = t, n().price_override = !0, c());
	}
	function d(e) {
		n().item_name = e.target.value, c();
	}
	function f(e) {
		n().unit = e.target.value, c();
	}
	function p(e) {
		n().category_type = e.target.value, c();
	}
	function m(e) {
		n().description = e.target.value || null, c();
	}
	function h() {
		R(s, !G(s));
	}
	function g() {
		t.ondelete?.(n().id);
	}
	var _ = pi(), v = Qt(_), y = B(v), b = B(y);
	Z(b, 21, () => r, Dr, (e, t) => {
		var n = di(), r = B(n, !0);
		M(n);
		var i = {};
		H(() => {
			Y(r, G(t)), i !== (i = G(t)) && (n.value = (n.__value = G(t)) ?? "");
		}), J(e, n);
	}), M(b);
	var x;
	Br(b), M(y);
	var S = V(y), C = B(S);
	Q(C), M(S);
	var w = V(S), T = B(w);
	Q(T), M(w);
	var E = V(w), D = B(E);
	Q(D), M(E);
	var ee = V(E), te = B(ee);
	Q(te), M(ee);
	var ne = V(ee), re = B(ne);
	M(ne);
	var ie = V(ne), ae = B(ie, !0);
	M(ie);
	var oe = V(ie), se = B(oe, !0);
	M(oe);
	var ce = V(oe), le = B(ce), ue = B(le), de = V(ue, 2);
	M(le), M(ce), M(v);
	var fe = V(v, 2), pe = (e) => {
		var t = fi(), r = B(t), i = B(r), a = V(B(i), 2);
		Q(a), M(i), M(r), M(t), H(() => Kr(a, n().description ?? "")), K("input", a, m), J(e, t);
	};
	X(fe, (e) => {
		G(s) && e(pe);
	}), H((e, t, r) => {
		x !== (x = n().category_type) && (b.value = (b.__value = n().category_type) ?? "", zr(b, n().category_type)), Kr(C, n().item_name), Kr(T, n().quantity), Kr(D, n().unit), Kr(te, e), Y(re, `${G(i) ?? ""}%`), Y(ae, t), Y(se, r), Rr(ue, 1, `p-0.5 transition-opacity
					${G(s) || n().description ? "opacity-100 text-[var(--color-sunburst)]" : "opacity-0 group-hover:opacity-100 text-white/20 hover:text-[var(--color-concrete)]"}`);
	}, [
		() => n().unit_price.toFixed(2),
		() => li(G(a)),
		() => li(G(o))
	]), K("change", b, p), K("input", C, d), K("input", T, l), K("input", D, f), K("input", te, u), K("click", ue, h), K("click", de, g), J(e, _), Fe();
}
mr([
	"change",
	"input",
	"click"
]);
//#endregion
//#region src/lib/Autocomplete.svelte
var hi = /* @__PURE__ */ q("<span class=\"text-xs text-[var(--color-sage)] font-[var(--font-ui)]\">Sub</span>"), gi = /* @__PURE__ */ q("<div class=\"flex items-center gap-2 text-xs\"><span class=\"font-mono text-[var(--color-sunburst)]\"> </span> <span class=\"text-white/30\"> </span></div>"), _i = /* @__PURE__ */ q("<button><div><span class=\"text-[var(--color-white)] font-[var(--font-body)]\"> </span> <span class=\"text-xs text-white/30 ml-2 font-[var(--font-body)]\"> </span></div> <!></button>"), vi = /* @__PURE__ */ q("<div class=\"absolute top-full left-0 right-0 mt-1 border border-white/[0.08] shadow-lg z-50 max-h-64 overflow-y-auto\" style=\"background: var(--color-granite);\"></div>"), yi = /* @__PURE__ */ q("<div class=\"absolute top-full left-0 right-0 mt-1 border border-white/[0.08] shadow-lg z-50\" style=\"background: var(--color-granite);\"><button class=\"w-full text-left px-3 py-2 text-sm text-[var(--color-muted-text)] hover:bg-white/[0.04] font-[var(--font-body)] transition-colors\">Add custom item: <span class=\"font-medium text-[var(--color-white)]\"> </span></button></div>"), bi = /* @__PURE__ */ q("<div class=\"relative\"><input type=\"text\" placeholder=\"Search items or type a name...\" class=\"w-full px-3 py-2 text-sm bg-white/[0.04] border border-white/[0.08] text-[var(--color-white)]\n			focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] outline-none placeholder-white/20 font-[var(--font-body)]\"/> <!></div>");
function xi(e, t) {
	Pe(t, !0);
	let n = $(t, "materialsDb", 19, () => []), r = $(t, "ratesDb", 19, () => []), i = $(t, "subcontractorsDb", 19, () => []), a = $(t, "categoryType", 3, "materials"), o = /* @__PURE__ */ L(""), s = /* @__PURE__ */ L(0), c, l = /* @__PURE__ */ I(() => {
		if (G(o).length < 1) return [];
		let e = G(o).toLowerCase(), t;
		return a() === "materials" ? t = n().map((e) => ({
			id: e.id,
			name: e.name,
			category: e.supplier || "",
			unit: e.unit,
			price: e.unit_price,
			type: "materials",
			source: "material"
		})) : (t = r().filter((e) => a() === "labor" ? e.category === "Labor" : a() === "equipment" ? e.category === "Equipment Rentals" : a() === "subs" ? e.category === "Subcontractors" : a() === "other" ? e.category === "Other" : !0).map((e) => ({
			id: e.id,
			name: e.name,
			category: e.category,
			unit: e.unit,
			price: e.rate,
			type: a(),
			source: "rate"
		})), a() === "subs" && i().length > 0 && (t = [...i().map((e) => ({
			id: e.id,
			name: e.company ? `${e.name} (${e.company})` : e.name,
			category: e.primary_trade || "Subcontractor",
			unit: "ls",
			price: 0,
			type: "subs",
			source: "subcontractor"
		})), ...t])), t.filter((t) => t.name.toLowerCase().includes(e) || t.category.toLowerCase().includes(e)).slice(0, 15);
	});
	function u(e) {
		e.key === "ArrowDown" ? (e.preventDefault(), R(s, Math.min(G(s) + 1, G(l).length - 1), !0)) : e.key === "ArrowUp" ? (e.preventDefault(), R(s, Math.max(G(s) - 1, 0), !0)) : e.key === "Enter" ? (e.preventDefault(), G(l).length > 0 && G(s) < G(l).length ? d(G(l)[G(s)]) : G(o).trim() && f()) : e.key === "Escape" && (e.preventDefault(), t.oncancel?.());
	}
	function d(e) {
		t.onselect?.({
			item_name: e.name,
			unit: e.unit,
			unit_price: e.price,
			category_type: e.type,
			is_custom: !1,
			material_id: e.source === "material" ? e.id : null,
			subcontractor_id: e.source === "subcontractor" ? e.id : null
		});
	}
	function f() {
		t.onselect?.({
			item_name: G(o).trim(),
			unit: "ea",
			unit_price: 0,
			category_type: a(),
			is_custom: !0,
			material_id: null,
			subcontractor_id: null
		});
	}
	mn(() => {
		R(s, 0);
	}), mn(() => {
		c?.focus();
	});
	var p = bi(), m = B(p);
	Q(m), ni(m, (e) => c = e, () => c);
	var h = V(m, 2), g = (e) => {
		var t = vi();
		Z(t, 21, () => G(l), Dr, (e, t, n) => {
			var r = _i(), i = B(r), a = B(i), o = B(a, !0);
			M(a);
			var c = V(a, 2), l = B(c, !0);
			M(c), M(i);
			var u = V(i, 2), f = (e) => {
				J(e, hi());
			}, p = (e) => {
				var n = gi(), r = B(n), i = B(r);
				M(r);
				var a = V(r, 2), o = B(a);
				M(a), M(n), H((e) => {
					Y(i, `$${e ?? ""}`), Y(o, `/ ${G(t).unit ?? ""}`);
				}, [() => G(t).price.toFixed(2)]), J(e, n);
			};
			X(u, (e) => {
				G(t).source === "subcontractor" ? e(f) : e(p, -1);
			}), M(r), H(() => {
				Rr(r, 1, `w-full text-left px-3 py-2 text-sm flex items-center justify-between
						hover:bg-white/[0.04] ${n === G(s) ? "bg-white/[0.06]" : ""} transition-colors`), Y(o, G(t).name), Y(l, G(t).category);
			}), K("click", r, () => d(G(t))), pr("mouseenter", r, () => R(s, n, !0)), J(e, r);
		}), M(t), J(e, t);
	}, _ = (e) => {
		var t = yi(), n = B(t), r = V(B(n)), i = B(r);
		M(r), M(n), M(t), H(() => Y(i, `"${G(o) ?? ""}"`)), K("click", n, f), J(e, t);
	};
	X(h, (e) => {
		G(l).length > 0 ? e(g) : G(o).length > 0 && e(_, 1);
	}), M(p), K("keydown", m, u), Qr(m, () => G(o), (e) => R(o, e)), J(e, p), Fe();
}
mr(["keydown", "click"]);
//#endregion
//#region node_modules/nanoid/url-alphabet/index.js
var Si = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict", Ci = (e = 21) => {
	let t = "", n = crypto.getRandomValues(new Uint8Array(e |= 0));
	for (; e--;) t += Si[n[e] & 63];
	return t;
}, wi = /* @__PURE__ */ q("<input type=\"text\" class=\"text-xs font-semibold text-[var(--color-concrete)] uppercase tracking-wide px-1 py-0.5 border border-white/[0.08] bg-white/[0.04] focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] font-[var(--font-ui)]\"/>"), Ti = /* @__PURE__ */ q("<span class=\"text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider cursor-pointer font-[var(--font-ui)]\" role=\"button\" tabindex=\"0\"> </span> <button class=\"opacity-0 group-hover/cg:opacity-100 text-white/15 hover:text-[var(--color-sunburst)] transition-opacity p-0.5\" title=\"Rename\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z\"></path></svg></button>", 1), Ei = /* @__PURE__ */ q("<tr><td colspan=\"9\" class=\"px-1 py-1\"><!></td></tr>"), Di = /* @__PURE__ */ q("<tr><td colspan=\"9\" class=\"px-1 py-0.5\"><button class=\"text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors ml-4\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Item</button></td></tr>"), Oi = /* @__PURE__ */ q("<!> <!>", 1), ki = /* @__PURE__ */ q("<tr class=\"border-b border-white/[0.04]\"><td colspan=\"9\" class=\"px-0 py-0\"><div class=\"flex items-center gap-2 py-1.5 px-1 bg-white/[0.02] group/vg\"><button class=\"flex items-center gap-1.5 text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider font-[var(--font-ui)]\"><svg fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> </button> <span class=\"text-xs text-white/20 font-[var(--font-body)]\"> </span> <button class=\"opacity-0 group-hover/vg:opacity-100 text-white/20 hover:text-red-400 transition-opacity p-0.5 ml-auto\" title=\"Ungroup items\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></button></div></td></tr> <!>", 1), Ai = /* @__PURE__ */ q("<tr class=\"border-b border-white/[0.04]\"><td colspan=\"9\" class=\"px-0 py-0\"><div class=\"flex items-center gap-2 py-1.5 px-1 bg-white/[0.02]\"><span class=\"flex items-center gap-1.5 text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider font-[var(--font-ui)]\"><svg class=\"w-2.5 h-2.5 text-white/25 rotate-90\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> </span></div></td></tr> <tr><td colspan=\"9\" class=\"px-1 py-1\"><!></td></tr>", 1), ji = /* @__PURE__ */ q("<table class=\"w-full\"><tbody><!><!><!></tbody></table>"), Mi = /* @__PURE__ */ q("<div class=\"mb-1\"><div class=\"flex items-center gap-2 py-1.5 px-1 bg-white/[0.02] border-b border-white/[0.04]\"><span class=\"flex items-center gap-1.5 text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider font-[var(--font-ui)]\"><svg class=\"w-2.5 h-2.5 text-white/25 rotate-90\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> </span></div> <!></div>"), Ni = /* @__PURE__ */ q("<div class=\"flex items-center gap-2 mt-1\"><input type=\"text\" placeholder=\"Group label\" class=\"flex-1 px-0 py-1 bg-transparent border-0 border-b border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/20\"/> <button class=\"px-2 py-1 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-xs font-[var(--font-ui)] font-bold uppercase tracking-wide hover:brightness-110\">Add</button> <button class=\"px-2 py-1 text-[var(--color-muted-text)] text-xs hover:text-[var(--color-white)] font-[var(--font-ui)]\">Cancel</button></div>"), Pi = /* @__PURE__ */ q("<button class=\"mt-1 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M4 6h16M4 12h16M4 18h7\"></path></svg> Add Label Group</button>"), Fi = /* @__PURE__ */ q("<div class=\"mt-1 mb-1\"><!></div>"), Ii = /* @__PURE__ */ q("<button class=\"mt-1 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Item</button>"), Li = /* @__PURE__ */ q("<div class=\"ml-5 mt-3 border-l-2 border-white/[0.06] pl-4\"><div class=\"flex items-center gap-2 mb-1 group/cg\"><!> <span class=\"text-xs text-white/25 font-[var(--font-body)]\"> </span> <button class=\"opacity-0 group-hover/cg:opacity-100 text-white/20 hover:text-red-400 transition-opacity p-0.5\" title=\"Delete section\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></button></div> <!> <!> <!> <!></div>");
function Ri(e, t) {
	Pe(t, !0);
	let n = $(t, "group", 7), r = /* @__PURE__ */ L(!1), i = /* @__PURE__ */ L(null), a = /* @__PURE__ */ L(!1), o = /* @__PURE__ */ L(z(n().name)), s = /* @__PURE__ */ L(!1), c = /* @__PURE__ */ L(""), l = /* @__PURE__ */ I(() => {
		let e = /* @__PURE__ */ new Map(), t = [];
		for (let r of n().line_items) r.visual_group ? (e.has(r.visual_group) || e.set(r.visual_group, []), e.get(r.visual_group).push(r)) : t.push(r);
		return {
			groups: e,
			ungrouped: t
		};
	}), u = z({});
	function d(e) {
		t.onsnapshot?.(), n().line_items.push({
			id: Ci(),
			category_type: e.category_type,
			item_name: e.item_name,
			quantity: 1,
			unit: e.unit,
			unit_price: e.unit_price,
			is_custom: e.is_custom,
			material_id: e.material_id,
			subcontractor_id: e.subcontractor_id || null,
			price_override: !1,
			description: null,
			sort_order: n().line_items.length,
			component_group_id: n().id,
			visual_group: G(i)
		}), R(r, !1), R(i, null), t.onchange?.();
	}
	function f(e = null) {
		R(i, e, !0), R(r, !0);
	}
	function p(e) {
		t.onsnapshot?.();
		let r = n().line_items.findIndex((t) => t.id === e);
		r !== -1 && (n().line_items.splice(r, 1), t.onchange?.());
	}
	function h() {
		R(o, n().name, !0), R(a, !0);
	}
	function g() {
		let e = G(o).trim();
		e && e !== n().name && (t.onsnapshot?.(), n().name = e, t.onchange?.()), R(a, !1);
	}
	function _() {
		let e = G(c).trim();
		e && (R(c, ""), R(s, !1), f(e));
	}
	function v(e) {
		t.onsnapshot?.();
		for (let t of n().line_items) t.visual_group === e && (t.visual_group = null);
		t.onchange?.();
	}
	var y = Li(), b = B(y), x = B(b), S = (e) => {
		var t = wi();
		Q(t), rn(t, !0), K("keydown", t, (e) => {
			e.key === "Enter" && g(), e.key === "Escape" && R(a, !1);
		}), pr("blur", t, g), Qr(t, () => G(o), (e) => R(o, e)), J(e, t);
	}, C = (e) => {
		var t = Ti(), r = Qt(t), i = B(r, !0);
		M(r);
		var a = V(r, 2);
		H(() => Y(i, n().name)), K("dblclick", r, h), K("click", a, h), J(e, t);
	};
	X(x, (e) => {
		G(a) ? e(S) : e(C, -1);
	});
	var w = V(x, 2), T = B(w);
	M(w);
	var E = V(w, 2);
	M(b);
	var D = V(b, 2), ee = (e) => {
		var n = ji(), a = B(n), o = B(a);
		Z(o, 17, () => G(l).ungrouped, (e) => e.id, (e, n) => {
			mi(e, {
				get item() {
					return G(n);
				},
				get globals() {
					return t.globals;
				},
				get markupOverrides() {
					return t.markupOverrides;
				},
				get markupEnabled() {
					return t.markupEnabled;
				},
				get onchange() {
					return t.onchange;
				},
				ondelete: p
			});
		});
		var s = V(o);
		Z(s, 17, () => [...G(l).groups.entries()], ([e, t]) => e, (e, n) => {
			var a = /* @__PURE__ */ I(() => m(G(n), 2));
			let o = () => G(a)[0], s = () => G(a)[1], c = /* @__PURE__ */ I(() => u[o()] ?? !1);
			var l = ki(), h = Qt(l), g = B(h), _ = B(g), y = B(_), b = B(y), x = V(b);
			M(y);
			var S = V(y, 2), C = B(S);
			M(S);
			var w = V(S, 2);
			M(_), M(g), M(h);
			var T = V(h, 2), E = (e) => {
				var n = Oi(), a = Qt(n);
				Z(a, 17, s, (e) => e.id, (e, n) => {
					mi(e, {
						get item() {
							return G(n);
						},
						get globals() {
							return t.globals;
						},
						get markupOverrides() {
							return t.markupOverrides;
						},
						get markupEnabled() {
							return t.markupEnabled;
						},
						get onchange() {
							return t.onchange;
						},
						ondelete: p
					});
				});
				var c = V(a, 2), l = (e) => {
					var n = Ei(), a = B(n);
					xi(B(a), {
						get materialsDb() {
							return t.materialsDb;
						},
						get ratesDb() {
							return t.ratesDb;
						},
						get subcontractorsDb() {
							return t.subcontractorsDb;
						},
						categoryType: "materials",
						onselect: d,
						oncancel: () => {
							R(r, !1), R(i, null);
						}
					}), M(a), M(n), J(e, n);
				}, u = (e) => {
					var t = Di(), n = B(t), r = B(n);
					M(n), M(t), K("click", r, () => f(o())), J(e, t);
				};
				X(c, (e) => {
					G(r) && G(i) === o() ? e(l) : e(u, -1);
				}), J(e, n);
			};
			X(T, (e) => {
				G(c) || e(E);
			}), H(() => {
				Rr(b, 0, `w-2.5 h-2.5 text-white/25 transition-transform ${G(c) ? "" : "rotate-90"}`), Y(x, ` ${o() ?? ""}`), Y(C, `(${s().length ?? ""})`);
			}), K("click", y, () => u[o()] = !G(c)), K("click", w, () => v(o())), J(e, l);
		});
		var c = V(s), h = (e) => {
			var n = Ai(), a = Qt(n), o = B(a), s = B(o), c = B(s), l = V(B(c));
			M(c), M(s), M(o), M(a);
			var u = V(a, 2), f = B(u);
			xi(B(f), {
				get materialsDb() {
					return t.materialsDb;
				},
				get ratesDb() {
					return t.ratesDb;
				},
				get subcontractorsDb() {
					return t.subcontractorsDb;
				},
				categoryType: "materials",
				onselect: d,
				oncancel: () => {
					R(r, !1), R(i, null);
				}
			}), M(f), M(u), H(() => Y(l, ` ${G(i) ?? ""}`)), J(e, n);
		}, g = /* @__PURE__ */ I(() => G(r) && G(i) !== null && !G(l).groups.has(G(i)));
		X(c, (e) => {
			G(g) && e(h);
		}), M(a), M(n), J(e, n);
	};
	X(D, (e) => {
		n().line_items.length > 0 && e(ee);
	});
	var te = V(D, 2), ne = (e) => {
		var n = Mi(), a = B(n), o = B(a), s = V(B(o));
		M(o), M(a), xi(V(a, 2), {
			get materialsDb() {
				return t.materialsDb;
			},
			get ratesDb() {
				return t.ratesDb;
			},
			get subcontractorsDb() {
				return t.subcontractorsDb;
			},
			categoryType: "materials",
			onselect: d,
			oncancel: () => {
				R(r, !1), R(i, null);
			}
		}), M(n), H(() => Y(s, ` ${G(i) ?? ""}`)), J(e, n);
	}, re = /* @__PURE__ */ I(() => G(r) && G(i) !== null && !G(l).groups.has(G(i)) && n().line_items.length === 0);
	X(te, (e) => {
		G(re) && e(ne);
	});
	var ie = V(te, 2), ae = (e) => {
		var t = Ni(), n = B(t);
		Q(n);
		var r = V(n, 2), i = V(r, 2);
		M(t), K("keydown", n, (e) => {
			e.key === "Enter" && _(), e.key === "Escape" && R(s, !1);
		}), Qr(n, () => G(c), (e) => R(c, e)), K("click", r, _), K("click", i, () => R(s, !1)), J(e, t);
	}, oe = (e) => {
		var t = Pi();
		K("click", t, () => R(s, !0)), J(e, t);
	};
	X(ie, (e) => {
		G(s) ? e(ae) : e(oe, -1);
	});
	var se = V(ie, 2), ce = (e) => {
		var n = Fi();
		xi(B(n), {
			get materialsDb() {
				return t.materialsDb;
			},
			get ratesDb() {
				return t.ratesDb;
			},
			get subcontractorsDb() {
				return t.subcontractorsDb;
			},
			categoryType: "materials",
			onselect: d,
			oncancel: () => {
				R(r, !1), R(i, null);
			}
		}), M(n), J(e, n);
	}, le = (e) => {
		var t = Ii();
		K("click", t, () => f(null)), J(e, t);
	};
	X(se, (e) => {
		G(r) && G(i) === null ? e(ce) : e(le, -1);
	}), M(y), H(() => Y(T, `(${n().line_items.length ?? ""})`)), K("click", E, () => t.ondelete?.(n().id)), J(e, y), Fe();
}
mr([
	"keydown",
	"dblclick",
	"click"
]);
//#endregion
//#region src/lib/SubcategoryBlock.svelte
var zi = /* @__PURE__ */ q("<input type=\"text\" class=\"px-2 py-0.5 border border-white/[0.08] text-sm font-medium text-[var(--color-white)] bg-white/[0.04] focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] font-[var(--font-ui)]\"/>"), Bi = /* @__PURE__ */ q("<span class=\"font-medium text-[var(--color-concrete)] text-sm font-[var(--font-ui)]\"> </span> <button class=\"opacity-0 group-hover/subcat:opacity-100 text-white/15 hover:text-[var(--color-sunburst)] transition-opacity p-0.5 shrink-0\" title=\"Rename\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z\"></path></svg></button>", 1), Vi = /* @__PURE__ */ q("<span class=\"text-xs px-1.5 py-0.5 bg-[var(--color-sunburst)]/10 text-[var(--color-sunburst)] font-medium font-[var(--font-ui)]\">overrides</span>"), Hi = /* @__PURE__ */ q("<span class=\"text-xs px-1.5 py-0.5 bg-[var(--color-sage)]/10 text-[var(--color-sage)] font-medium font-[var(--font-ui)]\"> </span>"), Ui = /* @__PURE__ */ q("<div class=\"text-center\"><span class=\"block text-xs font-medium text-[var(--color-concrete)] mb-1 font-[var(--font-ui)]\"> </span> <div class=\"flex items-center justify-center gap-1 mb-1\"><input type=\"checkbox\" class=\"w-3 h-3 border-white/[0.1] bg-white/[0.04] text-[var(--color-sunburst)] focus:ring-[var(--color-sunburst)] rounded-sm\"/> <span class=\"text-xs text-white/30 font-[var(--font-body)]\"> </span></div> <input type=\"number\" step=\"1\" min=\"0\"/> <div class=\"text-xs text-white/30 mt-0.5 font-mono\"> </div></div>"), Wi = /* @__PURE__ */ q("<div class=\"bg-white/[0.02] border border-white/[0.06] p-3 mb-3\"><div class=\"flex items-center justify-between mb-2\"><span class=\"text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider font-[var(--font-ui)]\">Markup Overrides</span> <button class=\"text-xs text-white/30 hover:text-[var(--color-white)] font-[var(--font-ui)]\">Close</button></div> <div class=\"grid grid-cols-5 gap-3\"></div> <div class=\"mt-3 pt-2 border-t border-white/[0.06] flex items-center gap-2\"><span class=\"text-xs font-medium text-[var(--color-concrete)] font-[var(--font-ui)]\">Lump Sum:</span> <span class=\"text-xs text-white/30\">$</span> <input type=\"number\" step=\"0.01\" min=\"0\" class=\"w-28 text-right text-xs font-mono px-2 py-1 border border-white/[0.08] bg-white/[0.04] text-[var(--color-white)]\n								focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)]\"/> <span class=\"text-xs text-white/30 font-[var(--font-body)]\">added post-markup</span></div></div>"), Gi = /* @__PURE__ */ q("<span> </span>"), Ki = /* @__PURE__ */ q("<div class=\"flex items-center gap-3 py-1.5 px-3 bg-[var(--color-sunburst)]/5 text-xs mb-2 cursor-pointer hover:bg-[var(--color-sunburst)]/8 transition-colors border-l-2 border-[var(--color-sunburst)]/30\" role=\"button\" tabindex=\"0\"><span class=\"font-medium text-[var(--color-sunburst)] font-[var(--font-ui)]\">Markup:</span> <!></div>"), qi = /* @__PURE__ */ q("<tr><td colspan=\"9\" class=\"px-1 py-1\"><!></td></tr>"), Ji = /* @__PURE__ */ q("<tr><td colspan=\"9\" class=\"px-1 py-0.5\"><button class=\"text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors ml-4\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Item</button></td></tr>"), Yi = /* @__PURE__ */ q("<!> <!>", 1), Xi = /* @__PURE__ */ q("<tr class=\"border-b border-white/[0.04]\"><td colspan=\"9\" class=\"px-0 py-0\"><div class=\"flex items-center gap-2 py-1.5 px-1 bg-white/[0.02] group/vg\"><button class=\"flex items-center gap-1.5 text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider font-[var(--font-ui)]\"><svg fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> </button> <span class=\"text-xs text-white/20 font-[var(--font-body)]\"> </span> <button class=\"opacity-0 group-hover/vg:opacity-100 text-white/20 hover:text-red-400 transition-opacity p-0.5 ml-auto\" title=\"Ungroup items\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></button></div></td></tr> <!>", 1), Zi = /* @__PURE__ */ q("<tr class=\"border-b border-white/[0.04]\"><td colspan=\"9\" class=\"px-0 py-0\"><div class=\"flex items-center gap-2 py-1.5 px-1 bg-white/[0.02]\"><span class=\"flex items-center gap-1.5 text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider font-[var(--font-ui)]\"><svg class=\"w-2.5 h-2.5 text-white/25 rotate-90\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> </span></div></td></tr> <tr><td colspan=\"9\" class=\"px-1 py-1\"><!></td></tr>", 1), Qi = /* @__PURE__ */ q("<table class=\"w-full\"><thead><tr class=\"text-xs text-white/30 uppercase tracking-wider border-b border-white/[0.06] font-[var(--font-ui)]\"><th class=\"px-1 py-1.5 text-left w-24\">Type</th><th class=\"px-1 py-1.5 text-left\">Name</th><th class=\"px-1 py-1.5 text-right w-20\">Qty</th><th class=\"px-1 py-1.5 text-center w-16\">Unit</th><th class=\"px-1 py-1.5 text-right w-24\">Price</th><th class=\"px-2 py-1.5 text-right w-16\">Markup</th><th class=\"px-2 py-1.5 text-right w-24\">w/ Markup</th><th class=\"px-2 py-1.5 text-right w-24\">Total</th><th class=\"w-16\"></th></tr></thead><tbody><!><!><!></tbody></table>"), $i = /* @__PURE__ */ q("<div class=\"mb-2\"><div class=\"flex items-center gap-2 py-1.5 px-1 bg-white/[0.02] border-b border-white/[0.04]\"><span class=\"flex items-center gap-1.5 text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider font-[var(--font-ui)]\"><svg class=\"w-2.5 h-2.5 text-white/25 rotate-90\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> </span></div> <!></div>"), ea = /* @__PURE__ */ q("<div class=\"flex items-center gap-2 ml-5 mt-2\"><input type=\"text\" placeholder=\"Section name\" class=\"flex-1 px-0 py-1 bg-transparent border-0 border-b border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/20\"/> <button class=\"px-2 py-1 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-xs font-[var(--font-ui)] font-bold uppercase tracking-wide hover:brightness-110\">Add</button> <button class=\"px-2 py-1 text-[var(--color-muted-text)] text-xs hover:text-[var(--color-white)] font-[var(--font-ui)]\">Cancel</button></div>"), ta = /* @__PURE__ */ q("<button class=\"ml-5 mt-2 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Section</button>"), na = /* @__PURE__ */ q("<div class=\"flex items-center gap-2 mt-2\"><input type=\"text\" placeholder=\"Group label (e.g. Interior, Phase 1)\" class=\"flex-1 px-0 py-1 bg-transparent border-0 border-b border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/20\"/> <button class=\"px-2 py-1 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-xs font-[var(--font-ui)] font-bold uppercase tracking-wide hover:brightness-110\">Add</button> <button class=\"px-2 py-1 text-[var(--color-muted-text)] text-xs hover:text-[var(--color-white)] font-[var(--font-ui)]\">Cancel</button></div>"), ra = /* @__PURE__ */ q("<button class=\"mt-1 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M4 6h16M4 12h16M4 18h7\"></path></svg> Add Label Group</button>"), ia = /* @__PURE__ */ q("<div class=\"mt-2\"><!></div>"), aa = /* @__PURE__ */ q("<button class=\"mt-2 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors\"><svg class=\"w-3.5 h-3.5\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Item</button>"), oa = /* @__PURE__ */ q("<span class=\"text-[var(--color-sage)]\"> </span>"), sa = /* @__PURE__ */ q("<div class=\"px-5 pb-3\"><!> <!> <!> <!> <!> <!> <!> <div class=\"flex justify-between items-center mt-3 pt-2 border-t border-white/[0.06] text-sm\"><span class=\"text-white/30 font-[var(--font-body)]\"> <!></span> <span class=\"font-mono font-semibold text-[var(--color-white)]\"> </span></div></div>"), ca = /* @__PURE__ */ q("<div class=\"border-t border-white/[0.06]\"><div class=\"flex items-center justify-between px-5 py-2 hover:bg-white/[0.02] transition-colors group/subcat\"><button class=\"flex items-center gap-2 text-left flex-1 min-w-0\"><svg fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> <!> <span class=\"text-xs text-white/25 font-[var(--font-body)]\"> </span> <!> <!></button> <div class=\"flex items-center gap-2\"><button class=\"text-xs px-1.5 py-0.5 bg-white/[0.04] text-[var(--color-muted-text)] hover:bg-white/[0.08] hover:text-[var(--color-concrete)] transition-colors\" title=\"Configure markup\"><svg class=\"w-3 h-3 inline\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4\"></path></svg></button> <span class=\"font-mono text-sm font-semibold text-[var(--color-white)]\"> </span> <button class=\"opacity-0 group-hover/subcat:opacity-100 text-white/20 hover:text-red-400 transition-opacity p-0.5\" title=\"Delete subcategory\"><svg class=\"w-3.5 h-3.5\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></button></div></div> <!></div>");
function la(e, t) {
	Pe(t, !0);
	let n = $(t, "subcat", 7), r = $(t, "collapsed", 3, !1), i = $(t, "materialsDb", 19, () => []), a = $(t, "ratesDb", 19, () => []), o = $(t, "subcontractorsDb", 19, () => []), s = /* @__PURE__ */ L(z(r())), c = /* @__PURE__ */ L(!1), l = /* @__PURE__ */ L(null), u = /* @__PURE__ */ L(!1), d = /* @__PURE__ */ L(z(n().name)), f = /* @__PURE__ */ L(!1), p = /* @__PURE__ */ L(""), h = /* @__PURE__ */ L(!1), g = /* @__PURE__ */ L(""), _ = /* @__PURE__ */ I(() => oi(n(), t.globals)), v = /* @__PURE__ */ I(() => [
		"materials",
		"labor",
		"equipment",
		"subs",
		"other"
	].map((e) => ({
		type: e,
		value: ri(e, t.globals, n().markup_overrides, n().markup_enabled),
		isOverride: n().markup_overrides[e] != null,
		isDisabled: !n().markup_enabled[e]
	}))), y = /* @__PURE__ */ L(!1), b = /* @__PURE__ */ I(() => G(v).some((e) => e.isOverride || e.isDisabled)), x = /* @__PURE__ */ I(() => n().line_items.length + n().component_groups.reduce((e, t) => e + t.line_items.length, 0)), S = /* @__PURE__ */ I(() => {
		let e = /* @__PURE__ */ new Map(), t = [];
		for (let r of n().line_items) r.visual_group ? (e.has(r.visual_group) || e.set(r.visual_group, []), e.get(r.visual_group).push(r)) : t.push(r);
		return {
			groups: e,
			ungrouped: t
		};
	}), C = z({});
	function w(e) {
		t.onsnapshot?.(), n().line_items.push({
			id: Ci(),
			category_type: e.category_type,
			item_name: e.item_name,
			quantity: 1,
			unit: e.unit,
			unit_price: e.unit_price,
			is_custom: e.is_custom,
			material_id: e.material_id,
			subcontractor_id: e.subcontractor_id || null,
			price_override: !1,
			description: null,
			sort_order: n().line_items.length,
			component_group_id: null,
			visual_group: G(l)
		}), R(c, !1), R(l, null), t.onchange?.();
	}
	function T(e = null) {
		R(l, e, !0), R(c, !0);
	}
	function E(e) {
		t.onsnapshot?.();
		let r = n().line_items.findIndex((t) => t.id === e);
		r !== -1 && (n().line_items.splice(r, 1), t.onchange?.());
	}
	function D() {
		R(d, n().name, !0), R(u, !0);
	}
	function ee() {
		let e = G(d).trim();
		e && e !== n().name && (t.onsnapshot?.(), n().name = e, t.onchange?.()), R(u, !1);
	}
	function te() {
		t.onsnapshot?.();
		let e = G(p).trim() || "New section";
		n().component_groups.push({
			id: Ci(),
			name: e,
			sort_order: n().component_groups.length,
			line_items: []
		}), R(p, ""), R(f, !1), t.onchange?.();
	}
	function ne(e) {
		t.onsnapshot?.();
		let r = n().component_groups.findIndex((t) => t.id === e);
		r !== -1 && (n().component_groups.splice(r, 1), t.onchange?.());
	}
	function re() {
		let e = G(g).trim();
		e && (R(g, ""), R(h, !1), T(e));
	}
	function ie(e) {
		t.onsnapshot?.();
		for (let t of n().line_items) t.visual_group === e && (t.visual_group = null);
		t.onchange?.();
	}
	function ae(e, r) {
		let i = r.target.value.trim();
		if (i === "") n().markup_overrides[e] = null;
		else {
			let t = parseFloat(i);
			!isNaN(t) && t >= 0 && (n().markup_overrides[e] = t);
		}
		t.onchange?.();
	}
	function oe(e) {
		n().markup_enabled[e] = !n().markup_enabled[e], t.onchange?.();
	}
	function se(e) {
		let r = parseFloat(e.target.value);
		!isNaN(r) && r >= 0 && (n().lump_sum = r, t.onchange?.());
	}
	let ce = {
		materials: "Mat",
		labor: "Lab",
		equipment: "Equip",
		subs: "Subs",
		other: "Other"
	};
	var le = ca(), ue = B(le), de = B(ue), fe = B(de), pe = V(fe, 2), me = (e) => {
		var t = zi();
		Q(t), rn(t, !0), K("click", t, (e) => e.stopPropagation()), K("keydown", t, (e) => {
			e.key === "Enter" && ee(), e.key === "Escape" && R(u, !1);
		}), pr("blur", t, ee), Qr(t, () => G(d), (e) => R(d, e)), J(e, t);
	}, he = (e) => {
		var t = Bi(), r = Qt(t), i = B(r, !0);
		M(r);
		var a = V(r, 2);
		H(() => Y(i, n().name)), K("click", a, (e) => {
			e.stopPropagation(), D();
		}), J(e, t);
	};
	X(pe, (e) => {
		G(u) ? e(me) : e(he, -1);
	});
	var ge = V(pe, 2), _e = B(ge);
	M(ge);
	var ve = V(ge, 2), O = (e) => {
		J(e, Vi());
	};
	X(ve, (e) => {
		G(b) && e(O);
	});
	var ye = V(ve, 2), be = (e) => {
		var t = Hi(), r = B(t);
		M(t), H((e) => Y(r, `+${e ?? ""} lump sum`), [() => li(n().lump_sum)]), J(e, t);
	};
	X(ye, (e) => {
		n().lump_sum > 0 && e(be);
	}), M(de);
	var xe = V(de, 2), Se = B(xe), k = V(Se, 2), Ce = B(k, !0);
	M(k);
	var A = V(k, 2);
	M(xe), M(ue);
	var j = V(ue, 2), we = (e) => {
		var r = sa(), s = B(r), u = (e) => {
			var r = Wi(), i = B(r), a = V(B(i), 2);
			M(i);
			var o = V(i, 2);
			Z(o, 20, () => [
				"materials",
				"labor",
				"equipment",
				"subs",
				"other"
			], Dr, (e, r) => {
				let i = /* @__PURE__ */ I(() => G(v).find((e) => e.type === r));
				var a = Ui(), o = B(a), s = B(o, !0);
				M(o);
				var c = V(o, 2), l = B(c);
				Q(l);
				var u = V(l, 2), d = B(u, !0);
				M(u), M(c);
				var f = V(c, 2);
				Q(f);
				var p = V(f, 2), m = B(p);
				M(p), M(a), H((e) => {
					Y(s, ce[r]), qr(l, n().markup_enabled[r]), Y(d, n().markup_enabled[r] ? "On" : "Off"), Kr(f, n().markup_overrides[r] ?? ""), Jr(f, "placeholder", `${t.globals[`${r}_markup`]}%`), f.disabled = !n().markup_enabled[r], Rr(f, 1, `w-full text-center text-xs font-mono px-1 py-1 border border-white/[0.08]
										focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)]
										${n().markup_enabled[r] ? "bg-white/[0.04] text-[var(--color-white)]" : "bg-white/[0.01] text-white/15"}
										${G(i)?.isOverride ? "border-[var(--color-sunburst)]/40 bg-[var(--color-sunburst)]/5" : ""}`), Y(m, `eff: ${e ?? ""}`);
				}, [() => ui(G(i)?.value ?? 0)]), K("change", l, () => oe(r)), K("input", f, (e) => ae(r, e)), J(e, a);
			}), M(o);
			var s = V(o, 2), c = V(B(s), 4);
			Q(c), Te(2), M(s), M(r), H(() => Kr(c, n().lump_sum)), K("click", a, () => R(y, !1)), K("input", c, se), J(e, r);
		}, d = (e) => {
			var t = Ki();
			Z(V(B(t), 2), 17, () => G(v), Dr, (e, t) => {
				var n = Gi(), r = B(n);
				M(n), H((e) => {
					Rr(n, 1, `${G(t).isDisabled ? "text-white/15 line-through" : G(t).isOverride ? "text-[var(--color-sunburst)]" : "text-white/30"} font-[var(--font-body)]`), Y(r, `${G(t).type ?? ""} ${e ?? ""}`);
				}, [() => ui(G(t).value)]), J(e, n);
			}), M(t), K("click", t, () => R(y, !0)), K("keydown", t, (e) => {
				(e.key === "Enter" || e.key === " ") && R(y, !0);
			}), J(e, t);
		};
		X(s, (e) => {
			G(y) ? e(u) : G(b) && e(d, 1);
		});
		var D = V(s, 2), ee = (e) => {
			var r = Qi(), s = V(B(r)), u = B(s);
			Z(u, 17, () => G(S).ungrouped, (e) => e.id, (e, r) => {
				mi(e, {
					get item() {
						return G(r);
					},
					get globals() {
						return t.globals;
					},
					get markupOverrides() {
						return n().markup_overrides;
					},
					get markupEnabled() {
						return n().markup_enabled;
					},
					get onchange() {
						return t.onchange;
					},
					ondelete: E
				});
			});
			var d = V(u);
			Z(d, 17, () => [...G(S).groups.entries()], ([e, t]) => e, (e, r) => {
				var s = /* @__PURE__ */ I(() => m(G(r), 2));
				let u = () => G(s)[0], d = () => G(s)[1], f = /* @__PURE__ */ I(() => C[u()] ?? !1);
				var p = Xi(), h = Qt(p), g = B(h), _ = B(g), v = B(_), y = B(v), b = V(y);
				M(v);
				var x = V(v, 2), S = B(x);
				M(x);
				var D = V(x, 2);
				M(_), M(g), M(h);
				var ee = V(h, 2), te = (e) => {
					var r = Yi(), s = Qt(r);
					Z(s, 17, d, (e) => e.id, (e, r) => {
						mi(e, {
							get item() {
								return G(r);
							},
							get globals() {
								return t.globals;
							},
							get markupOverrides() {
								return n().markup_overrides;
							},
							get markupEnabled() {
								return n().markup_enabled;
							},
							get onchange() {
								return t.onchange;
							},
							ondelete: E
						});
					});
					var f = V(s, 2), p = (e) => {
						var t = qi(), n = B(t);
						xi(B(n), {
							get materialsDb() {
								return i();
							},
							get ratesDb() {
								return a();
							},
							get subcontractorsDb() {
								return o();
							},
							categoryType: "materials",
							onselect: w,
							oncancel: () => {
								R(c, !1), R(l, null);
							}
						}), M(n), M(t), J(e, t);
					}, m = (e) => {
						var t = Ji(), n = B(t), r = B(n);
						M(n), M(t), K("click", r, () => T(u())), J(e, t);
					};
					X(f, (e) => {
						G(c) && G(l) === u() ? e(p) : e(m, -1);
					}), J(e, r);
				};
				X(ee, (e) => {
					G(f) || e(te);
				}), H(() => {
					Rr(y, 0, `w-2.5 h-2.5 text-white/25 transition-transform ${G(f) ? "" : "rotate-90"}`), Y(b, ` ${u() ?? ""}`), Y(S, `(${d().length ?? ""})`);
				}), K("click", v, () => C[u()] = !G(f)), K("click", D, () => ie(u())), J(e, p);
			});
			var f = V(d), p = (e) => {
				var t = Zi(), n = Qt(t), r = B(n), s = B(r), u = B(s), d = V(B(u));
				M(u), M(s), M(r), M(n);
				var f = V(n, 2), p = B(f);
				xi(B(p), {
					get materialsDb() {
						return i();
					},
					get ratesDb() {
						return a();
					},
					get subcontractorsDb() {
						return o();
					},
					categoryType: "materials",
					onselect: w,
					oncancel: () => {
						R(c, !1), R(l, null);
					}
				}), M(p), M(f), H(() => Y(d, ` ${G(l) ?? ""}`)), J(e, t);
			}, h = /* @__PURE__ */ I(() => G(c) && G(l) !== null && !G(S).groups.has(G(l)));
			X(f, (e) => {
				G(h) && e(p);
			}), M(s), M(r), J(e, r);
		};
		X(D, (e) => {
			(G(x) > 0 || n().line_items.length > 0) && e(ee);
		});
		var le = V(D, 2), ue = (e) => {
			var t = $i(), n = B(t), r = B(n), s = V(B(r));
			M(r), M(n), xi(V(n, 2), {
				get materialsDb() {
					return i();
				},
				get ratesDb() {
					return a();
				},
				get subcontractorsDb() {
					return o();
				},
				categoryType: "materials",
				onselect: w,
				oncancel: () => {
					R(c, !1), R(l, null);
				}
			}), M(t), H(() => Y(s, ` ${G(l) ?? ""}`)), J(e, t);
		}, de = /* @__PURE__ */ I(() => G(c) && G(l) !== null && !G(S).groups.has(G(l)) && G(x) === 0);
		X(le, (e) => {
			G(de) && e(ue);
		});
		var fe = V(le, 2);
		Z(fe, 17, () => n().component_groups, (e) => e.id, (e, r) => {
			Ri(e, {
				get group() {
					return G(r);
				},
				get globals() {
					return t.globals;
				},
				get markupOverrides() {
					return n().markup_overrides;
				},
				get markupEnabled() {
					return n().markup_enabled;
				},
				get onchange() {
					return t.onchange;
				},
				get onsnapshot() {
					return t.onsnapshot;
				},
				ondelete: ne,
				get materialsDb() {
					return i();
				},
				get ratesDb() {
					return a();
				},
				get subcontractorsDb() {
					return o();
				}
			});
		});
		var pe = V(fe, 2), me = (e) => {
			var t = ea(), n = B(t);
			Q(n);
			var r = V(n, 2), i = V(r, 2);
			M(t), K("keydown", n, (e) => {
				e.key === "Enter" && te(), e.key === "Escape" && R(f, !1);
			}), Qr(n, () => G(p), (e) => R(p, e)), K("click", r, te), K("click", i, () => R(f, !1)), J(e, t);
		}, he = (e) => {
			var t = ta();
			K("click", t, () => R(f, !0)), J(e, t);
		};
		X(pe, (e) => {
			G(f) ? e(me) : e(he, -1);
		});
		var ge = V(pe, 2), _e = (e) => {
			var t = na(), n = B(t);
			Q(n);
			var r = V(n, 2), i = V(r, 2);
			M(t), K("keydown", n, (e) => {
				e.key === "Enter" && re(), e.key === "Escape" && R(h, !1);
			}), Qr(n, () => G(g), (e) => R(g, e)), K("click", r, re), K("click", i, () => R(h, !1)), J(e, t);
		}, ve = (e) => {
			var t = ra();
			K("click", t, () => R(h, !0)), J(e, t);
		};
		X(ge, (e) => {
			G(h) ? e(_e) : e(ve, -1);
		});
		var O = V(ge, 2), ye = (e) => {
			var t = ia();
			xi(B(t), {
				get materialsDb() {
					return i();
				},
				get ratesDb() {
					return a();
				},
				get subcontractorsDb() {
					return o();
				},
				categoryType: "materials",
				onselect: w,
				oncancel: () => {
					R(c, !1), R(l, null);
				}
			}), M(t), J(e, t);
		}, be = (e) => {
			var t = aa();
			K("click", t, () => T(null)), J(e, t);
		};
		X(O, (e) => {
			G(c) && G(l) === null ? e(ye) : e(be, -1);
		});
		var xe = V(O, 2), Se = B(xe), k = B(Se), Ce = V(k), A = (e) => {
			var t = oa(), r = B(t);
			M(t), H((e) => Y(r, `+ ${e ?? ""} lump sum`), [() => li(n().lump_sum)]), J(e, t);
		};
		X(Ce, (e) => {
			n().lump_sum > 0 && e(A);
		}), M(Se);
		var j = V(Se, 2), we = B(j, !0);
		M(j), M(xe), M(r), H((e, t) => {
			Y(k, `Subtotal: ${e ?? ""} `), Y(we, t);
		}, [() => li(G(_).base), () => li(G(_).withMarkup)]), J(e, r);
	};
	X(j, (e) => {
		G(s) || e(we);
	}), M(le), H((e) => {
		Rr(fe, 0, `w-3 h-3 text-white/25 transition-transform shrink-0 ${G(s) ? "" : "rotate-90"}`), Y(_e, `${G(x) ?? ""} item${G(x) === 1 ? "" : "s"}`), Y(Ce, e);
	}, [() => li(G(_).withMarkup)]), K("click", de, () => R(s, !G(s))), K("click", Se, () => R(y, !G(y))), K("click", A, () => t.ondelete?.(n().id)), J(e, le), Fe();
}
mr([
	"click",
	"keydown",
	"change",
	"input"
]);
//#endregion
//#region src/lib/SectionBlock.svelte
var ua = /* @__PURE__ */ q("<input type=\"text\" class=\"bg-white/[0.06] text-[var(--color-white)] px-2 py-0.5 text-sm font-semibold border border-white/[0.08] focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] font-[var(--font-ui)] uppercase tracking-wide\"/>"), da = /* @__PURE__ */ q("<span class=\"font-semibold font-[var(--font-ui)] uppercase tracking-wide text-[var(--color-white)]\"> </span> <button class=\"opacity-0 group-hover/sec:opacity-100 text-white/20 hover:text-[var(--color-sunburst)] transition-opacity p-0.5 shrink-0\" title=\"Rename category\"><svg class=\"w-3.5 h-3.5\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z\"></path></svg></button>", 1), fa = /* @__PURE__ */ q("<div class=\"px-5 py-10 text-center text-[var(--color-muted-text)] text-sm font-[var(--font-body)]\">No subcategories yet</div>"), pa = /* @__PURE__ */ q("<div class=\"flex items-center gap-2 mt-2\"><input type=\"text\" placeholder=\"Subcategory name\" class=\"flex-1 px-0 py-1.5 bg-transparent border-0 border-b border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/20\"/> <button class=\"px-3 py-1.5 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-xs font-[var(--font-ui)] font-bold uppercase tracking-wide hover:brightness-110\">Add</button> <button class=\"px-2 py-1.5 text-[var(--color-muted-text)] text-xs hover:text-[var(--color-white)] font-[var(--font-ui)]\">Cancel</button></div>"), ma = /* @__PURE__ */ q("<button class=\"mt-2 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] uppercase tracking-wide transition-colors\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Subcategory</button>"), ha = /* @__PURE__ */ q("<!> <div class=\"px-5 pb-4\"><!></div>", 1), ga = /* @__PURE__ */ q("<div class=\"mb-6 border border-white/[0.06] overflow-hidden\" style=\"background: var(--color-ink);\"><div class=\"flex items-center justify-between px-5 py-3 group/sec\" style=\"background: var(--color-granite);\"><button class=\"flex items-center gap-3 hover:bg-white/[0.04] -ml-2 px-2 py-1 transition-colors text-left flex-1 min-w-0\"><svg fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> <!> <span class=\"text-xs text-white/30 font-[var(--font-body)]\"> </span></button> <div class=\"flex items-center gap-3\"><span class=\"font-mono font-semibold text-[var(--color-white)]\"> </span> <button class=\"text-white/20 hover:text-red-400 transition-colors p-1\" title=\"Delete category\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16\"></path></svg></button></div></div> <!></div>");
function _a(e, t) {
	Pe(t, !0);
	let n = $(t, "section", 7), r = $(t, "collapsed", 3, !1), i = $(t, "materialsDb", 19, () => []), a = $(t, "ratesDb", 19, () => []), o = $(t, "subcontractorsDb", 19, () => []), s = /* @__PURE__ */ L(z(r())), c = /* @__PURE__ */ L(!1), l = /* @__PURE__ */ L(z(n().name)), u = /* @__PURE__ */ L(!1), d = /* @__PURE__ */ L(""), f = /* @__PURE__ */ I(() => si(n(), t.globals)), p = /* @__PURE__ */ I(() => n().subcategories.reduce((e, t) => e + t.line_items.length + t.component_groups.reduce((e, t) => e + t.line_items.length, 0), 0));
	function m() {
		R(l, n().name, !0), R(c, !0);
	}
	function h() {
		let e = G(l).trim();
		e && e !== n().name && (t.onsnapshot?.(), n().name = e, t.onchange?.()), R(c, !1);
	}
	function g() {
		t.onsnapshot?.();
		let e = G(d).trim() || "New subcategory";
		n().subcategories.push({
			id: Ci(),
			name: e,
			sort_order: n().subcategories.length,
			lump_sum: 0,
			markup_overrides: {
				materials: null,
				labor: null,
				equipment: null,
				subs: null,
				other: null
			},
			markup_enabled: {
				materials: !0,
				labor: !0,
				equipment: !0,
				subs: !0,
				other: !0
			},
			line_items: [],
			component_groups: []
		}), R(d, ""), R(u, !1), t.onchange?.();
	}
	function _(e) {
		t.onsnapshot?.();
		let r = n().subcategories.findIndex((t) => t.id === e);
		r !== -1 && (n().subcategories.splice(r, 1), t.onchange?.());
	}
	var v = ga(), y = B(v), b = B(y), x = B(b), S = V(x, 2), C = (e) => {
		var t = ua();
		Q(t), rn(t, !0), K("click", t, (e) => e.stopPropagation()), K("keydown", t, (e) => {
			e.key === "Enter" && h(), e.key === "Escape" && R(c, !1);
		}), pr("blur", t, h), Qr(t, () => G(l), (e) => R(l, e)), J(e, t);
	}, w = (e) => {
		var t = da(), r = Qt(t), i = B(r, !0);
		M(r);
		var a = V(r, 2);
		H(() => Y(i, n().name)), K("click", a, (e) => {
			e.stopPropagation(), m();
		}), J(e, t);
	};
	X(S, (e) => {
		G(c) ? e(C) : e(w, -1);
	});
	var T = V(S, 2), E = B(T);
	M(T), M(b);
	var D = V(b, 2), ee = B(D), te = B(ee, !0);
	M(ee);
	var ne = V(ee, 2);
	M(D), M(y);
	var re = V(y, 2), ie = (e) => {
		var r = ha(), s = Qt(r), c = (e) => {
			J(e, fa());
		}, l = (e) => {
			var r = xr();
			Z(Qt(r), 17, () => n().subcategories, (e) => e.id, (e, n) => {
				la(e, {
					get subcat() {
						return G(n);
					},
					get globals() {
						return t.globals;
					},
					get onchange() {
						return t.onchange;
					},
					get onsnapshot() {
						return t.onsnapshot;
					},
					ondelete: _,
					get materialsDb() {
						return i();
					},
					get ratesDb() {
						return a();
					},
					get subcontractorsDb() {
						return o();
					}
				});
			}), J(e, r);
		};
		X(s, (e) => {
			n().subcategories.length === 0 && !G(u) ? e(c) : e(l, -1);
		});
		var f = V(s, 2), p = B(f), m = (e) => {
			var t = pa(), n = B(t);
			Q(n);
			var r = V(n, 2), i = V(r, 2);
			M(t), K("keydown", n, (e) => {
				e.key === "Enter" && g(), e.key === "Escape" && R(u, !1);
			}), Qr(n, () => G(d), (e) => R(d, e)), K("click", r, g), K("click", i, () => R(u, !1)), J(e, t);
		}, h = (e) => {
			var t = ma();
			K("click", t, () => R(u, !0)), J(e, t);
		};
		X(p, (e) => {
			G(u) ? e(m) : e(h, -1);
		}), M(f), J(e, r);
	};
	X(re, (e) => {
		G(s) || e(ie);
	}), M(v), H((e) => {
		Rr(x, 0, `w-4 h-4 text-white/30 transition-transform shrink-0 ${G(s) ? "" : "rotate-90"}`), Y(E, `${n().subcategories.length ?? ""} subcategor${n().subcategories.length === 1 ? "y" : "ies"}
				· ${G(p) ?? ""} item${G(p) === 1 ? "" : "s"}`), Y(te, e);
	}, [() => li(G(f).withMarkup)]), K("click", b, () => R(s, !G(s))), K("click", ne, () => t.ondelete?.(n().id)), J(e, v), Fe();
}
mr(["click", "keydown"]);
//#endregion
//#region src/lib/FooterSummary.svelte
var va = /* @__PURE__ */ q("<div class=\"flex flex-col\"><span class=\"text-xs text-white/30 uppercase tracking-wider font-[var(--font-ui)]\"> </span> <span class=\"font-mono text-[var(--color-concrete)]\"> </span></div>"), ya = /* @__PURE__ */ q("<span class=\"text-[var(--color-muted-text)] text-xs font-[var(--font-body)]\">No items yet</span>"), ba = /* @__PURE__ */ q("<div class=\"fixed bottom-0 right-0 border-t-2 z-20\" style=\"left: var(--sidebar-width, 220px); background: var(--color-granite); border-color: var(--color-sunburst);\"><div class=\"flex items-center justify-between px-5 py-3\"><div class=\"flex items-center gap-6 text-sm\"><div class=\"flex flex-col\"><span class=\"text-xs text-white/30 uppercase tracking-wider font-[var(--font-ui)]\">Base Cost</span> <span class=\"font-mono text-[var(--color-concrete)]\"> </span></div> <div class=\"w-px h-8 bg-white/[0.08]\"></div> <!> <!></div> <div class=\"flex flex-col items-end\"><span class=\"text-xs text-white/30 uppercase tracking-wider font-[var(--font-ui)]\">Total</span> <span class=\"font-mono text-lg font-bold text-[var(--color-sunburst)]\"> </span></div></div></div>");
function xa(e, t) {
	Pe(t, !0);
	let n = /* @__PURE__ */ I(() => ci(t.estimate)), r = /* @__PURE__ */ I(() => Object.entries(G(n).byType).filter(([, e]) => e > 0).map(([e, t]) => ({
		type: e,
		value: t
	}))), i = {
		materials: "Materials",
		labor: "Labor",
		equipment: "Equipment",
		subs: "Subs",
		other: "Other"
	};
	var a = ba(), o = B(a), s = B(o), c = B(s), l = V(B(c), 2), u = B(l, !0);
	M(l), M(c);
	var d = V(c, 4);
	Z(d, 17, () => G(r), Dr, (e, t) => {
		let n = () => G(t).type, r = () => G(t).value;
		var a = va(), o = B(a), s = B(o, !0);
		M(o);
		var c = V(o, 2), l = B(c, !0);
		M(c), M(a), H((e) => {
			Y(s, i[n()]), Y(l, e);
		}, [() => li(r())]), J(e, a);
	});
	var f = V(d, 2), p = (e) => {
		J(e, ya());
	};
	X(f, (e) => {
		G(r).length === 0 && e(p);
	}), M(s);
	var m = V(s, 2), h = V(B(m), 2), g = B(h, !0);
	M(h), M(m), M(o), M(a), H((e, t) => {
		Y(u, e), Y(g, t);
	}, [() => li(G(n).base), () => li(G(n).withMarkup)]), J(e, a), Fe();
}
//#endregion
//#region src/lib/SaveStatus.svelte
var Sa = /* @__PURE__ */ q("<button class=\"ml-1 px-1.5 py-0.5 rounded bg-[var(--color-sunburst)]/20 hover:bg-[var(--color-sunburst)]/30 text-[var(--color-sunburst)] transition-colors\" title=\"Save now (Ctrl+S)\">Save</button>"), Ca = /* @__PURE__ */ q("<button class=\"ml-1 px-1.5 py-0.5 rounded bg-red-500/20 hover:bg-red-500/30 text-red-300 transition-colors\" title=\"Retry save\">Retry</button>"), wa = /* @__PURE__ */ q("<div><span></span> <span> </span> <!> <!></div>");
function Ta(e, t) {
	Pe(t, !0);
	let n = /* @__PURE__ */ L(z(Date.now()));
	mn(() => {
		let e = setInterval(() => {
			R(n, Date.now(), !0);
		}, 1e4);
		return () => clearInterval(e);
	});
	let r = /* @__PURE__ */ I(() => {
		switch (t.status) {
			case "clean":
				if (t.savedAt) {
					let e = Math.round((G(n) - t.savedAt.getTime()) / 1e3);
					return e < 5 ? "Saved just now" : e < 60 ? `Saved ${e}s ago` : `Saved ${Math.round(e / 60)}m ago`;
				}
				return "Up to date";
			case "dirty": return "Unsaved changes";
			case "saving": return "Saving...";
			case "error": return "Save failed";
			default: return "";
		}
	}), i = /* @__PURE__ */ I(() => {
		switch (t.status) {
			case "clean": return "text-[var(--color-sage)]";
			case "dirty": return "text-[var(--color-sunburst)]";
			case "saving": return "text-blue-400";
			case "error": return "text-red-400";
			default: return "text-white/40";
		}
	}), a = /* @__PURE__ */ I(() => {
		switch (t.status) {
			case "clean": return "bg-[var(--color-sage)]";
			case "dirty": return "bg-[var(--color-sunburst)]";
			case "saving": return "bg-blue-400 animate-pulse";
			case "error": return "bg-red-400 animate-pulse";
			default: return "bg-white/40";
		}
	});
	var o = wa(), s = B(o), c = V(s, 2), l = B(c, !0);
	M(c);
	var u = V(c, 2), d = (e) => {
		var n = Sa();
		K("click", n, function(...e) {
			t.onsave?.apply(this, e);
		}), J(e, n);
	};
	X(u, (e) => {
		t.status === "dirty" && t.onsave && e(d);
	});
	var f = V(u, 2), p = (e) => {
		var n = Ca();
		K("click", n, function(...e) {
			t.onsave?.apply(this, e);
		}), J(e, n);
	};
	X(f, (e) => {
		t.status === "error" && t.onsave && e(p);
	}), M(o), H(() => {
		Rr(o, 1, `flex items-center gap-1.5 text-xs font-[var(--font-ui)] ${G(i) ?? ""}`), Rr(s, 1, `w-2 h-2 rounded-full ${G(a) ?? ""}`), Y(l, G(r));
	}), J(e, o), Fe();
}
mr(["click"]);
//#endregion
//#region src/lib/autosave.svelte.js
function Ea(e, t = 2e3) {
	let n = /* @__PURE__ */ L("clean"), r = /* @__PURE__ */ L(null), i = null, a = null;
	function o(e) {
		a = e;
	}
	function s() {
		R(n, "dirty"), i && clearTimeout(i), i = setTimeout(() => c(), t);
	}
	async function c() {
		if (i &&= (clearTimeout(i), null), !a) return null;
		let o = a();
		if (!o) return null;
		R(n, "saving");
		try {
			let t = await fetch(`/api/estimate/${e}`, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(o)
			});
			if (!t.ok) {
				let e = await t.json().catch(() => ({ error: `HTTP ${t.status}` }));
				throw Error(e.error || `Save failed: ${t.status}`);
			}
			let i = await t.json();
			return R(n, "clean"), R(r, /* @__PURE__ */ new Date(), !0), i;
		} catch (e) {
			return R(n, "error"), console.error("Auto-save failed:", e.message), i = setTimeout(() => c(), t * 2), null;
		}
	}
	function l() {
		i && clearTimeout(i);
	}
	return {
		register: o,
		markDirty: s,
		save: c,
		destroy: l,
		get status() {
			return G(n);
		},
		get savedAt() {
			return G(r);
		}
	};
}
//#endregion
//#region src/lib/undo.svelte.js
var Da = 20;
function Oa() {
	let e = z([]), t = /* @__PURE__ */ I(() => e.length > 0);
	function n(t) {
		let n = JSON.parse(JSON.stringify({
			globals: t.globals,
			sections: t.sections
		}));
		e.push(n), e.length > Da && e.shift();
	}
	function r(t) {
		if (e.length === 0) return !1;
		let n = e.pop();
		return t.globals = n.globals, t.sections = n.sections, !0;
	}
	function i() {
		e.length = 0;
	}
	return {
		get canUndo() {
			return G(t);
		},
		get depth() {
			return e.length;
		},
		snapshot: n,
		undo: r,
		clear: i
	};
}
//#endregion
//#region src/EstimateBuilder.svelte
var ka = /* @__PURE__ */ q("<div class=\"flex items-center justify-center h-64\"><div class=\"text-[var(--color-muted-text)] font-[var(--font-body)]\">Loading estimate...</div></div>"), Aa = /* @__PURE__ */ q("<div class=\"flex items-center justify-center h-64\"><div class=\"text-red-400 font-[var(--font-body)]\"> </div></div>"), ja = /* @__PURE__ */ q("<label class=\"flex items-center gap-1 text-xs text-[var(--color-concrete)] font-[var(--font-ui)]\"><span> </span> <input type=\"number\" step=\"1\" min=\"0\" class=\"w-10 text-right bg-transparent border-0 border-b border-white/[0.08] p-0 font-mono text-xs text-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none\"/> <span class=\"text-white/30\">%</span></label>"), Ma = /* @__PURE__ */ q("<div class=\"text-center py-16 text-[var(--color-muted-text)]\"><p class=\"text-lg font-[var(--font-ui)]\">No categories yet</p> <p class=\"text-sm mt-1 font-[var(--font-body)]\">Add a category to start building your estimate.</p> <button class=\"mt-4 px-4 py-2 bg-[var(--color-granite)] border border-white/[0.06] text-[var(--color-sunburst)] text-sm font-[var(--font-ui)] font-semibold uppercase tracking-wide hover:border-[var(--color-sunburst)] transition-colors\">+ Add First Category</button></div>"), Na = /* @__PURE__ */ q("<div class=\"flex items-center gap-2 mt-4\"><input type=\"text\" placeholder=\"Category name (e.g. Framing, Roofing)\" class=\"flex-1 px-0 py-2 bg-transparent border-0 border-b-2 border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/20\"/> <button class=\"px-4 py-2 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-xs font-[var(--font-ui)] font-bold uppercase tracking-wide hover:brightness-110\">Add</button> <button class=\"px-3 py-2 text-[var(--color-muted-text)] text-xs hover:text-[var(--color-white)] font-[var(--font-ui)]\">Cancel</button></div>"), Pa = /* @__PURE__ */ q("<button class=\"mt-6 w-full py-3 border border-dashed border-white/[0.08] text-sm text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] hover:border-[var(--color-sunburst)]/30 flex items-center justify-center gap-2 font-[var(--font-ui)] uppercase tracking-wide transition-colors\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Category</button>"), Fa = /* @__PURE__ */ q("<div class=\"estimate-builder pb-20\"><div class=\"sticky top-0 z-10 border-b border-white/[0.06]\" style=\"background: var(--color-granite);\"><div class=\"flex items-center justify-between px-5 py-2\"><div class=\"flex items-center gap-4\"><button class=\"text-xs px-2 py-1 bg-white/[0.06] hover:bg-white/[0.1] disabled:opacity-30 disabled:cursor-not-allowed transition-colors flex items-center gap-1 font-[var(--font-ui)] text-[var(--color-concrete)]\" title=\"Undo (Ctrl+Z)\"><svg class=\"w-3.5 h-3.5\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M3 10h10a5 5 0 015 5v2M3 10l4-4M3 10l4 4\"></path></svg> Undo</button> <!></div> <div class=\"flex items-center gap-3 text-sm\"><span class=\"text-[var(--color-sage)] font-[var(--font-ui)] uppercase tracking-wider text-xs font-semibold\">Markup</span> <!></div></div></div> <div class=\"px-5 pt-6\"><!> <!></div> <!></div>");
function Ia(e, t) {
	Pe(t, !0);
	let n = /* @__PURE__ */ L(null), r = /* @__PURE__ */ L(null), i = /* @__PURE__ */ L(!0), a = /* @__PURE__ */ L(!1), o = /* @__PURE__ */ L(""), s = Ea(t.projectId), c = Oa();
	async function l() {
		try {
			let e = await fetch(`/api/estimate/${t.projectId}`);
			if (!e.ok) throw Error(`Failed to load estimate: ${e.status}`);
			R(n, await e.json(), !0), s.register(() => G(n));
		} catch (e) {
			R(r, e.message, !0);
		} finally {
			R(i, !1);
		}
	}
	mn(() => {
		t.projectId && l();
	}), mn(() => {
		let e = document.getElementById("estimate-root");
		e && (e.dataset.dirty = s.status === "dirty" || s.status === "saving" ? "true" : "false");
	}), mn(() => {
		function e(e) {
			(s.status === "dirty" || s.status === "saving") && e.preventDefault();
		}
		function t(e) {
			if (s.status !== "dirty" && s.status !== "saving") return;
			let t = e.target.closest("a[href]");
			t && t.origin === window.location.origin && (t.closest("#estimate-root") || (e.preventDefault(), confirm("You have unsaved changes. Leave without saving?") && (window.location.href = t.href)));
		}
		function r(e) {
			(e.ctrlKey || e.metaKey) && e.key === "z" && !e.shiftKey && (e.preventDefault(), G(n) && c.undo(G(n)) && u()), (e.ctrlKey || e.metaKey) && e.key === "s" && (e.preventDefault(), s.status === "dirty" && s.save());
		}
		return window.addEventListener("beforeunload", e), document.addEventListener("click", t, !0), window.addEventListener("keydown", r), () => {
			window.removeEventListener("beforeunload", e), document.removeEventListener("click", t, !0), window.removeEventListener("keydown", r), s.destroy();
		};
	});
	function u() {
		s.markDirty();
	}
	function d() {
		G(n) && c.snapshot(G(n));
	}
	function f() {
		d();
		let e = G(o).trim() || "New Section";
		G(n).sections.push({
			id: Ci(),
			name: e,
			sort_order: G(n).sections.length,
			subcategories: []
		}), R(o, ""), R(a, !1), u();
	}
	function p(e) {
		d();
		let t = G(n).sections.findIndex((t) => t.id === e);
		t !== -1 && (G(n).sections.splice(t, 1), u());
	}
	function m(e, t) {
		let r = parseFloat(t.target.value);
		!isNaN(r) && r >= 0 && (G(n).globals[`${e}_markup`] = r, u());
	}
	let h = [
		{
			key: "materials",
			label: "Materials"
		},
		{
			key: "labor",
			label: "Labor"
		},
		{
			key: "equipment",
			label: "Equipment"
		},
		{
			key: "subs",
			label: "Subs"
		},
		{
			key: "other",
			label: "Other"
		}
	];
	var g = xr(), _ = Qt(g), v = (e) => {
		J(e, ka());
	}, y = (e) => {
		var t = Aa(), n = B(t), i = B(n, !0);
		M(n), M(t), H(() => Y(i, G(r))), J(e, t);
	}, b = (e) => {
		var t = Fa(), r = B(t), i = B(r), l = B(i), g = B(l);
		Ta(V(g, 2), {
			get status() {
				return s.status;
			},
			get savedAt() {
				return s.savedAt;
			},
			onsave: () => s.save()
		}), M(l);
		var _ = V(l, 2);
		Z(V(B(_), 2), 17, () => h, Dr, (e, t) => {
			var r = ja(), i = B(r), a = B(i, !0);
			M(i);
			var o = V(i, 2);
			Q(o), Te(2), M(r), H(() => {
				Y(a, G(t).label), Kr(o, G(n).globals[`${G(t).key}_markup`]);
			}), K("input", o, (e) => m(G(t).key, e)), J(e, r);
		}), M(_), M(i), M(r);
		var v = V(r, 2), y = B(v), b = (e) => {
			var t = Ma(), n = V(B(t), 4);
			M(t), K("click", n, () => R(a, !0)), J(e, t);
		}, x = (e) => {
			var t = xr();
			Z(Qt(t), 17, () => G(n).sections, (e) => e.id, (e, t) => {
				_a(e, {
					get section() {
						return G(t);
					},
					get globals() {
						return G(n).globals;
					},
					onchange: u,
					onsnapshot: d,
					ondelete: p,
					get materialsDb() {
						return G(n).materials_db;
					},
					get ratesDb() {
						return G(n).rates_db;
					},
					get subcontractorsDb() {
						return G(n).subcontractors_db;
					}
				});
			}), J(e, t);
		};
		X(y, (e) => {
			G(n).sections.length === 0 && !G(a) ? e(b) : e(x, -1);
		});
		var S = V(y, 2), C = (e) => {
			var t = Na(), n = B(t);
			Q(n);
			var r = V(n, 2), i = V(r, 2);
			M(t), K("keydown", n, (e) => {
				e.key === "Enter" && f(), e.key === "Escape" && R(a, !1);
			}), Qr(n, () => G(o), (e) => R(o, e)), K("click", r, f), K("click", i, () => R(a, !1)), J(e, t);
		}, w = (e) => {
			var t = Pa();
			K("click", t, () => R(a, !0)), J(e, t);
		};
		X(S, (e) => {
			G(a) ? e(C) : G(n).sections.length > 0 && e(w, 1);
		}), M(v), xa(V(v, 2), { get estimate() {
			return G(n);
		} }), M(t), H(() => g.disabled = !c.canUndo), K("click", g, () => {
			c.undo(G(n)) && u();
		}), J(e, t);
	};
	X(_, (e) => {
		G(i) ? e(v) : G(r) ? e(y, 1) : G(n) && e(b, 2);
	}), J(e, g), Fe();
}
mr([
	"click",
	"input",
	"keydown"
]);
//#endregion
//#region src/main.js
var La = document.getElementById("estimate-root");
if (La) {
	let e = La.dataset.projectId;
	Sr(Ia, {
		target: La,
		props: { projectId: e }
	});
}
//#endregion
