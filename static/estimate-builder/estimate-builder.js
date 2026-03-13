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
var m = 1024, h = 2048, g = 4096, _ = 8192, v = 16384, y = 32768, b = 1 << 25, x = 65536, S = 1 << 19, ee = 1 << 20, te = 1 << 25, C = 65536, ne = 1 << 21, re = 1 << 22, ie = 1 << 23, ae = Symbol("$state"), oe = Symbol("legacy props"), se = Symbol(""), w = new class extends Error {
	name = "StaleReactionError";
	message = "The reaction that called `getAbortSignal()` was re-run or destroyed";
}(), ce = !!globalThis.document?.contentType && /* @__PURE__ */ globalThis.document.contentType.includes("xml");
//#endregion
//#region node_modules/svelte/src/internal/client/errors.js
function le() {
	throw Error("https://svelte.dev/e/async_derived_orphan");
}
function ue(e, t, n) {
	throw Error("https://svelte.dev/e/each_key_duplicate");
}
function de(e) {
	throw Error("https://svelte.dev/e/effect_in_teardown");
}
function fe() {
	throw Error("https://svelte.dev/e/effect_in_unowned_derived");
}
function pe(e) {
	throw Error("https://svelte.dev/e/effect_orphan");
}
function me() {
	throw Error("https://svelte.dev/e/effect_update_depth_exceeded");
}
function he(e) {
	throw Error("https://svelte.dev/e/props_invalid_value");
}
function ge() {
	throw Error("https://svelte.dev/e/state_descriptors_fixed");
}
function _e() {
	throw Error("https://svelte.dev/e/state_prototype_fixed");
}
function ve() {
	throw Error("https://svelte.dev/e/state_unsafe_mutation");
}
function ye() {
	throw Error("https://svelte.dev/e/svelte_boundary_reset_onerror");
}
//#endregion
//#region node_modules/svelte/src/constants.js
var be = {}, T = Symbol(), xe = "http://www.w3.org/1999/xhtml";
function Se(e) {
	console.warn("https://svelte.dev/e/hydration_mismatch");
}
function Ce() {
	console.warn("https://svelte.dev/e/select_multiple_invalid_value");
}
function we() {
	console.warn("https://svelte.dev/e/svelte_boundary_reset_noop");
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/hydration.js
var E = !1;
function Te(e) {
	E = e;
}
var D;
function Ee(e) {
	if (e === null) throw Se(), be;
	return D = e;
}
function De() {
	return Ee(/* @__PURE__ */ $t(D));
}
function O(e) {
	if (E) {
		if (/* @__PURE__ */ $t(D) !== null) throw Se(), be;
		D = e;
	}
}
function Oe(e = 1) {
	if (E) {
		for (var t = e, n = D; t--;) n = /* @__PURE__ */ $t(n);
		D = n;
	}
}
function ke(e = !0) {
	for (var t = 0, n = D;;) {
		if (n.nodeType === 8) {
			var r = n.data;
			if (r === "]") {
				if (t === 0) return n;
				--t;
			} else (r === "[" || r === "[!" || r[0] === "[" && !isNaN(Number(r.slice(1)))) && (t += 1);
		}
		var i = /* @__PURE__ */ $t(n);
		e && n.remove(), n = i;
	}
}
function Ae(e) {
	if (!e || e.nodeType !== 8) throw Se(), be;
	return e.data;
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/equality.js
function je(e) {
	return e === this.v;
}
function Me(e, t) {
	return e == e ? e !== t || typeof e == "object" && !!e || typeof e == "function" : t == t;
}
function Ne(e) {
	return !Me(e, this.v);
}
//#endregion
//#region node_modules/svelte/src/internal/flags/index.js
var Pe = !1, Fe = !1, k = null;
function Ie(e) {
	k = e;
}
function Le(e, t = !1, n) {
	k = {
		p: k,
		i: !1,
		c: null,
		e: null,
		s: e,
		x: null,
		r: W,
		l: Fe && !t ? {
			s: null,
			u: null,
			$: []
		} : null
	};
}
function Re(e) {
	var t = k, n = t.e;
	if (n !== null) {
		t.e = null;
		for (var r of n) gn(r);
	}
	return e !== void 0 && (t.x = e), t.i = !0, k = t.p, e ?? {};
}
function ze() {
	return !Fe || k !== null && k.l === null;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/task.js
var Be = [];
function Ve() {
	var e = Be;
	Be = [], f(e);
}
function He(e) {
	if (Be.length === 0 && !nt) {
		var t = Be;
		queueMicrotask(() => {
			t === Be && Ve();
		});
	}
	Be.push(e);
}
function Ue() {
	for (; Be.length > 0;) Ve();
}
function We(e) {
	var t = W;
	if (t === null) return H.f |= ie, e;
	if (!(t.f & 32768) && !(t.f & 4)) throw e;
	Ge(e, t);
}
function Ge(e, t) {
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
var Ke = ~(h | g | m);
function A(e, t) {
	e.f = e.f & Ke | t;
}
function qe(e) {
	e.f & 512 || e.deps === null ? A(e, m) : A(e, g);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/utils.js
function Je(e) {
	if (e !== null) for (let t of e) !(t.f & 2) || !(t.f & 65536) || (t.f ^= C, Je(t.deps));
}
function Ye(e, t, n) {
	e.f & 2048 ? t.add(e) : e.f & 4096 && n.add(e), Je(e.deps), A(e, m);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/store.js
var Xe = !1, Ze = !1;
function Qe(e) {
	var t = Ze;
	try {
		return Ze = !1, [e(), Ze];
	} finally {
		Ze = t;
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/batch.js
var $e = /* @__PURE__ */ new Set(), j = null, et = null, M = null, tt = null, nt = !1, rt = !1, it = null, at = null, ot = 0, st = 1, ct = class e {
	id = st++;
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
			for (var n of t.d) A(n, h), this.schedule(n);
			for (n of t.m) A(n, g), this.schedule(n);
		}
	}
	#d() {
		ot++ > 1e3 && ut();
		let t = this.#a;
		this.#a = [], this.apply();
		var n = it = [], r = [], i = at = [];
		for (let e of t) try {
			this.#f(e, n, r);
		} catch (t) {
			throw _t(e), t;
		}
		if (j = null, i.length > 0) {
			var a = e.ensure();
			for (let e of i) a.schedule(e);
		}
		if (it = null, at = null, this.#u()) {
			this.#p(r), this.#p(n);
			for (let [e, t] of this.#c) gt(e, t);
		} else {
			this.#n === 0 && $e.delete(this), this.#o.clear(), this.#s.clear();
			for (let e of this.#e) e(this);
			this.#e.clear(), et = this, ft(r), ft(n), et = null, this.#i?.resolve();
		}
		var o = j;
		if (this.#a.length > 0) {
			let e = o ??= this;
			e.#a.push(...this.#a.filter((t) => !e.#a.includes(t)));
		}
		o !== null && ($e.add(o), o.#d()), $e.has(this) || this.#m();
	}
	#f(e, t, n) {
		e.f ^= m;
		for (var r = e.first; r !== null;) {
			var i = r.f, a = (i & 96) != 0;
			if (!(a && i & 1024 || i & 8192 || this.#c.has(r)) && r.fn !== null) {
				a ? r.f ^= m : i & 4 ? t.push(r) : Pe && i & 16777224 ? n.push(r) : Kn(r) && (i & 16 && this.#s.add(r), Zn(r));
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
		for (var t = 0; t < e.length; t += 1) Ye(e[t], this.#o, this.#s);
	}
	capture(e, t) {
		t !== T && !this.previous.has(e) && this.previous.set(e, t), e.f & 8388608 || (this.current.set(e, e.v), M?.set(e, e.v));
	}
	activate() {
		j = this;
	}
	deactivate() {
		j = null, M = null;
	}
	flush() {
		try {
			if (rt = !0, j = this, !this.#u()) {
				for (let e of this.#o) this.#s.delete(e), A(e, h), this.schedule(e);
				for (let e of this.#s) A(e, g), this.schedule(e);
			}
			this.#d();
		} finally {
			ot = 0, tt = null, it = null, at = null, rt = !1, j = null, M = null, It.clear();
		}
	}
	discard() {
		for (let e of this.#t) e(this);
		this.#t.clear();
	}
	#m() {
		for (let s of $e) {
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
					for (var a of t) pt(a, n, r, i);
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
		--this.#n, e && --this.#r, !(this.#l || t) && (this.#l = !0, He(() => {
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
		if (j === null) {
			let t = j = new e();
			rt || ($e.add(j), nt || He(() => {
				j === t && t.flush();
			}));
		}
		return j;
	}
	apply() {
		if (!Pe || !this.is_fork && $e.size === 1) {
			M = null;
			return;
		}
		M = new Map(this.current);
		for (let e of $e) if (e !== this) for (let [t, n] of e.previous) M.has(t) || M.set(t, n);
	}
	schedule(e) {
		if (tt = e, e.b?.is_pending && e.f & 16777228 && !(e.f & 32768)) {
			e.b.defer_effect(e);
			return;
		}
		for (var t = e; t.parent !== null;) {
			t = t.parent;
			var n = t.f;
			if (it !== null && t === W && (Pe || (H === null || !(H.f & 2)) && !Xe)) return;
			if (n & 96) {
				if (!(n & 1024)) return;
				t.f ^= m;
			}
		}
		this.#a.push(t);
	}
};
function lt(e) {
	var t = nt;
	nt = !0;
	try {
		var n;
		for (e && (j !== null && !j.is_fork && j.flush(), n = e());;) {
			if (Ue(), j === null) return n;
			j.flush();
		}
	} finally {
		nt = t;
	}
}
function ut() {
	try {
		me();
	} catch (e) {
		Ge(e, tt);
	}
}
var dt = null;
function ft(e) {
	var t = e.length;
	if (t !== 0) {
		for (var n = 0; n < t;) {
			var r = e[n++];
			if (!(r.f & 24576) && Kn(r) && (dt = /* @__PURE__ */ new Set(), Zn(r), r.deps === null && r.first === null && r.nodes === null && r.teardown === null && r.ac === null && En(r), dt?.size > 0)) {
				It.clear();
				for (let e of dt) {
					if (e.f & 24576) continue;
					let t = [e], n = e.parent;
					for (; n !== null;) dt.has(n) && (dt.delete(n), t.push(n)), n = n.parent;
					for (let e = t.length - 1; e >= 0; e--) {
						let n = t[e];
						n.f & 24576 || Zn(n);
					}
				}
				dt.clear();
			}
		}
		dt = null;
	}
}
function pt(e, t, n, r) {
	if (!n.has(e) && (n.add(e), e.reactions !== null)) for (let i of e.reactions) {
		let e = i.f;
		e & 2 ? pt(i, t, n, r) : e & 4194320 && !(e & 2048) && mt(i, t, r) && (A(i, h), ht(i));
	}
}
function mt(e, t, r) {
	let i = r.get(e);
	if (i !== void 0) return i;
	if (e.deps !== null) for (let i of e.deps) {
		if (n.call(t, i)) return !0;
		if (i.f & 2 && mt(i, t, r)) return r.set(i, !0), !0;
	}
	return r.set(e, !1), !1;
}
function ht(e) {
	j.schedule(e);
}
function gt(e, t) {
	if (!(e.f & 32 && e.f & 1024)) {
		e.f & 2048 ? t.d.push(e) : e.f & 4096 && t.m.push(e), A(e, m);
		for (var n = e.first; n !== null;) gt(n, t), n = n.next;
	}
}
function _t(e) {
	A(e, m);
	for (var t = e.first; t !== null;) _t(t), t = t.next;
}
//#endregion
//#region node_modules/svelte/src/reactivity/create-subscriber.js
function vt(e) {
	let t = 0, n = Rt(0), r;
	return () => {
		pn() && (J(n), bn(() => (t === 0 && (r = tr(() => e(() => Ht(n)))), t += 1, () => {
			He(() => {
				--t, t === 0 && (r?.(), r = void 0, Ht(n));
			});
		})));
	};
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/boundary.js
var yt = x | S;
function bt(e, t, n, r) {
	new xt(e, t, n, r);
}
var xt = class {
	parent;
	is_pending = !1;
	transform_error;
	#e;
	#t = E ? D : null;
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
	#h = vt(() => (this.#m = Rt(this.#l), () => {
		this.#m = null;
	}));
	constructor(e, t, n, r) {
		this.#e = e, this.#n = t, this.#r = (e) => {
			var t = W;
			t.b = this, t.f |= 128, n(e);
		}, this.parent = W.b, this.transform_error = r ?? this.parent?.transform_error ?? ((e) => e), this.#i = xn(() => {
			if (E) {
				let e = this.#t;
				De();
				let t = e.data === "[!";
				if (e.data.startsWith("[?")) {
					let t = JSON.parse(e.data.slice(2));
					this.#_(t);
				} else t ? this.#v() : this.#g();
			} else this.#y();
		}, yt), E && (this.#e = D);
	}
	#g() {
		try {
			this.#a = B(() => this.#r(this.#e));
		} catch (e) {
			this.error(e);
		}
	}
	#_(e) {
		let t = this.#n.failed;
		t && (this.#s = B(() => {
			t(this.#e, () => e, () => () => {});
		}));
	}
	#v() {
		let e = this.#n.pending;
		e && (this.is_pending = !0, this.#o = B(() => e(this.#e)), He(() => {
			var e = this.#c = document.createDocumentFragment(), t = I();
			e.append(t), this.#a = this.#x(() => B(() => this.#r(t))), this.#u === 0 && (this.#e.before(e), this.#c = null, Dn(this.#o, () => {
				this.#o = null;
			}), this.#b(j));
		}));
	}
	#y() {
		try {
			if (this.is_pending = this.has_pending_snippet(), this.#u = 0, this.#l = 0, this.#a = B(() => {
				this.#r(this.#e);
			}), this.#u > 0) {
				var e = this.#c = document.createDocumentFragment();
				jn(this.#a, e);
				let t = this.#n.pending;
				this.#o = B(() => t(this.#e));
			} else this.#b(j);
		} catch (e) {
			this.error(e);
		}
	}
	#b(e) {
		this.is_pending = !1;
		for (let t of this.#f) A(t, h), e.schedule(t);
		for (let t of this.#p) A(t, g), e.schedule(t);
		this.#f.clear(), this.#p.clear();
	}
	defer_effect(e) {
		Ye(e, this.#f, this.#p);
	}
	is_rendered() {
		return !this.is_pending && (!this.parent || this.parent.is_rendered());
	}
	has_pending_snippet() {
		return !!this.#n.pending;
	}
	#x(e) {
		var t = W, n = H, r = k;
		Ln(this.#i), U(this.#i), Ie(this.#i.ctx);
		try {
			return ct.ensure(), e();
		} catch (e) {
			return We(e), null;
		} finally {
			Ln(t), U(n), Ie(r);
		}
	}
	#S(e, t) {
		if (!this.has_pending_snippet()) {
			this.parent && this.parent.#S(e, t);
			return;
		}
		this.#u += e, this.#u === 0 && (this.#b(t), this.#o && Dn(this.#o, () => {
			this.#o = null;
		}), this.#c &&= (this.#e.before(this.#c), null));
	}
	update_pending_count(e, t) {
		this.#S(e, t), this.#l += e, !(!this.#m || this.#d) && (this.#d = !0, He(() => {
			this.#d = !1, this.#m && Bt(this.#m, this.#l);
		}));
	}
	get_effect_pending() {
		return this.#h(), J(this.#m);
	}
	error(e) {
		var t = this.#n.onerror;
		let n = this.#n.failed;
		if (!t && !n) throw e;
		this.#a &&= (V(this.#a), null), this.#o &&= (V(this.#o), null), this.#s &&= (V(this.#s), null), E && (Ee(this.#t), Oe(), Ee(ke()));
		var r = !1, i = !1;
		let a = () => {
			if (r) {
				we();
				return;
			}
			r = !0, i && ye(), this.#s !== null && Dn(this.#s, () => {
				this.#s = null;
			}), this.#x(() => {
				this.#y();
			});
		}, o = (e) => {
			try {
				i = !0, t?.(e, a), i = !1;
			} catch (e) {
				Ge(e, this.#i && this.#i.parent);
			}
			n && (this.#s = this.#x(() => {
				try {
					return B(() => {
						var t = W;
						t.b = this, t.f |= 128, n(this.#e, () => e, () => a);
					});
				} catch (e) {
					return Ge(e, this.#i.parent), null;
				}
			}));
		};
		He(() => {
			var t;
			try {
				t = this.transform_error(e);
			} catch (e) {
				Ge(e, this.#i && this.#i.parent);
				return;
			}
			typeof t == "object" && t && typeof t.then == "function" ? t.then(o, (e) => Ge(e, this.#i && this.#i.parent)) : o(t);
		});
	}
};
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/async.js
function St(e, t, n, r) {
	let i = ze() ? Et : Ot;
	var a = e.filter((e) => !e.settled);
	if (n.length === 0 && a.length === 0) {
		r(t.map(i));
		return;
	}
	var o = W, s = Ct(), c = a.length === 1 ? a[0].promise : a.length > 1 ? Promise.all(a.map((e) => e.promise)) : null;
	function l(e) {
		s();
		try {
			r(e);
		} catch (e) {
			o.f & 16384 || Ge(e, o);
		}
		wt();
	}
	if (n.length === 0) {
		c.then(() => l(t.map(i)));
		return;
	}
	var u = Tt();
	function d() {
		Promise.all(n.map((e) => /* @__PURE__ */ Dt(e))).then((e) => l([...t.map(i), ...e])).catch((e) => Ge(e, o)).finally(() => u());
	}
	c ? c.then(() => {
		s(), d(), wt();
	}) : d();
}
function Ct() {
	var e = W, t = H, n = k, r = j;
	return function(i = !0) {
		Ln(e), U(t), Ie(n), i && !(e.f & 16384) && (r?.activate(), r?.apply());
	};
}
function wt(e = !0) {
	Ln(null), U(null), Ie(null), e && j?.deactivate();
}
function Tt() {
	var e = W.b, t = j, n = e.is_rendered();
	return e.update_pending_count(1, t), t.increment(n), (r = !1) => {
		e.update_pending_count(-1, t), t.decrement(n, r);
	};
}
/* @__NO_SIDE_EFFECTS__ */
function Et(e) {
	var t = 2 | h, n = H !== null && H.f & 2 ? H : null;
	return W !== null && (W.f |= S), {
		ctx: k,
		deps: null,
		effects: null,
		equals: je,
		f: t,
		fn: e,
		reactions: null,
		rv: 0,
		v: T,
		wv: 0,
		parent: n ?? W,
		ac: null
	};
}
/* @__NO_SIDE_EFFECTS__ */
function Dt(e, t, n) {
	let r = W;
	r === null && le();
	var i = void 0, a = Rt(T), o = !H, s = /* @__PURE__ */ new Map();
	return yn(() => {
		var t = W, n = p();
		i = n.promise;
		try {
			Promise.resolve(e()).then(n.resolve, n.reject).finally(wt);
		} catch (e) {
			n.reject(e), wt();
		}
		var c = j;
		if (o) {
			if (t.f & 32768) var l = Tt();
			if (r.b.is_rendered()) s.get(c)?.reject(w), s.delete(c);
			else {
				for (let e of s.values()) e.reject(w);
				s.clear();
			}
			s.set(c, n);
		}
		let u = (e, n = void 0) => {
			if (l && l(n === w), !(n === w || t.f & 16384)) {
				if (c.activate(), n) a.f |= ie, Bt(a, n);
				else {
					a.f & 8388608 && (a.f ^= ie), Bt(a, e);
					for (let [e, t] of s) {
						if (s.delete(e), e === c) break;
						t.reject(w);
					}
				}
				c.deactivate();
			}
		};
		n.promise.then(u, (e) => u(null, e || "unknown"));
	}), mn(() => {
		for (let e of s.values()) e.reject(w);
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
function N(e) {
	let t = /* @__PURE__ */ Et(e);
	return Pe || Rn(t), t;
}
/* @__NO_SIDE_EFFECTS__ */
function Ot(e) {
	let t = /* @__PURE__ */ Et(e);
	return t.equals = Ne, t;
}
function kt(e) {
	var t = e.effects;
	if (t !== null) {
		e.effects = null;
		for (var n = 0; n < t.length; n += 1) V(t[n]);
	}
}
function At(e) {
	for (var t = e.parent; t !== null;) {
		if (!(t.f & 2)) return t.f & 16384 ? null : t;
		t = t.parent;
	}
	return null;
}
function jt(e) {
	var t, n = W;
	Ln(At(e));
	try {
		e.f &= ~C, kt(e), t = Jn(e);
	} finally {
		Ln(n);
	}
	return t;
}
function Mt(e) {
	var t = jt(e);
	if (!e.equals(t) && (e.wv = Gn(), (!j?.is_fork || e.deps === null) && (e.v = t, e.deps === null))) {
		A(e, m);
		return;
	}
	Pn || (M === null ? qe(e) : (pn() || j?.is_fork) && M.set(e, t));
}
function Nt(e) {
	if (e.effects !== null) for (let t of e.effects) (t.teardown || t.ac) && (t.teardown?.(), t.ac?.abort(w), t.teardown = d, t.ac = null, Xn(t, 0), Cn(t));
}
function Pt(e) {
	if (e.effects !== null) for (let t of e.effects) t.teardown && Zn(t);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/sources.js
var Ft = /* @__PURE__ */ new Set(), It = /* @__PURE__ */ new Map(), Lt = !1;
function Rt(e, t) {
	return {
		f: 0,
		v: e,
		reactions: null,
		equals: je,
		rv: 0,
		wv: 0
	};
}
/* @__NO_SIDE_EFFECTS__ */
function P(e, t) {
	let n = Rt(e, t);
	return Rn(n), n;
}
/* @__NO_SIDE_EFFECTS__ */
function zt(e, t = !1, n = !0) {
	let r = Rt(e);
	return t || (r.equals = Ne), Fe && n && k !== null && k.l !== null && (k.l.s ??= []).push(r), r;
}
function F(e, t, r = !1) {
	return H !== null && (!In || H.f & 131072) && ze() && H.f & 4325394 && (G === null || !n.call(G, e)) && ve(), Bt(e, r ? Wt(t) : t, at);
}
function Bt(e, t, n = null) {
	if (!e.equals(t)) {
		var r = e.v;
		Pn ? It.set(e, t) : It.set(e, r), e.v = t;
		var i = ct.ensure();
		if (i.capture(e, r), e.f & 2) {
			let t = e;
			e.f & 2048 && jt(t), qe(t);
		}
		e.wv = Gn(), Ut(e, h, n), ze() && W !== null && W.f & 1024 && !(W.f & 96) && (zn === null ? Bn([e]) : zn.push(e)), !i.is_fork && Ft.size > 0 && !Lt && Vt();
	}
	return t;
}
function Vt() {
	Lt = !1;
	for (let e of Ft) e.f & 1024 && A(e, g), Kn(e) && Zn(e);
	Ft.clear();
}
function Ht(e) {
	F(e, e.v + 1);
}
function Ut(e, t, n) {
	var r = e.reactions;
	if (r !== null) for (var i = ze(), a = r.length, o = 0; o < a; o++) {
		var s = r[o], c = s.f;
		if (!(!i && s === W)) {
			var l = (c & h) === 0;
			if (l && A(s, t), c & 2) {
				var u = s;
				M?.delete(u), c & 65536 || (c & 512 && (s.f |= C), Ut(u, g, n));
			} else if (l) {
				var d = s;
				c & 16 && dt !== null && dt.add(d), n === null ? ht(d) : n.push(d);
			}
		}
	}
}
function Wt(t) {
	if (typeof t != "object" || !t || ae in t) return t;
	let n = l(t);
	if (n !== s && n !== c) return t;
	var r = /* @__PURE__ */ new Map(), i = e(t), o = /* @__PURE__ */ P(0), u = null, d = Un, f = (e) => {
		if (Un === d) return e();
		var t = H, n = Un;
		U(null), Wn(d);
		var r = e();
		return U(t), Wn(n), r;
	};
	return i && r.set("length", /* @__PURE__ */ P(t.length, u)), new Proxy(t, {
		defineProperty(e, t, n) {
			(!("value" in n) || n.configurable === !1 || n.enumerable === !1 || n.writable === !1) && ge();
			var i = r.get(t);
			return i === void 0 ? f(() => {
				var e = /* @__PURE__ */ P(n.value, u);
				return r.set(t, e), e;
			}) : F(i, n.value, !0), !0;
		},
		deleteProperty(e, t) {
			var n = r.get(t);
			if (n === void 0) {
				if (t in e) {
					let e = f(() => /* @__PURE__ */ P(T, u));
					r.set(t, e), Ht(o);
				}
			} else F(n, T), Ht(o);
			return !0;
		},
		get(e, n, i) {
			if (n === ae) return t;
			var o = r.get(n), s = n in e;
			if (o === void 0 && (!s || a(e, n)?.writable) && (o = f(() => /* @__PURE__ */ P(Wt(s ? e[n] : T), u)), r.set(n, o)), o !== void 0) {
				var c = J(o);
				return c === T ? void 0 : c;
			}
			return Reflect.get(e, n, i);
		},
		getOwnPropertyDescriptor(e, t) {
			var n = Reflect.getOwnPropertyDescriptor(e, t);
			if (n && "value" in n) {
				var i = r.get(t);
				i && (n.value = J(i));
			} else if (n === void 0) {
				var a = r.get(t), o = a?.v;
				if (a !== void 0 && o !== T) return {
					enumerable: !0,
					configurable: !0,
					value: o,
					writable: !0
				};
			}
			return n;
		},
		has(e, t) {
			if (t === ae) return !0;
			var n = r.get(t), i = n !== void 0 && n.v !== T || Reflect.has(e, t);
			return (n !== void 0 || W !== null && (!i || a(e, t)?.writable)) && (n === void 0 && (n = f(() => /* @__PURE__ */ P(i ? Wt(e[t]) : T, u)), r.set(t, n)), J(n) === T) ? !1 : i;
		},
		set(e, t, n, s) {
			var c = r.get(t), l = t in e;
			if (i && t === "length") for (var d = n; d < c.v; d += 1) {
				var p = r.get(d + "");
				p === void 0 ? d in e && (p = f(() => /* @__PURE__ */ P(T, u)), r.set(d + "", p)) : F(p, T);
			}
			if (c === void 0) (!l || a(e, t)?.writable) && (c = f(() => /* @__PURE__ */ P(void 0, u)), F(c, Wt(n)), r.set(t, c));
			else {
				l = c.v !== T;
				var m = f(() => Wt(n));
				F(c, m);
			}
			var h = Reflect.getOwnPropertyDescriptor(e, t);
			if (h?.set && h.set.call(s, n), !l) {
				if (i && typeof t == "string") {
					var g = r.get("length"), _ = Number(t);
					Number.isInteger(_) && _ >= g.v && F(g, _ + 1);
				}
				Ht(o);
			}
			return !0;
		},
		ownKeys(e) {
			J(o);
			var t = Reflect.ownKeys(e).filter((e) => {
				var t = r.get(e);
				return t === void 0 || t.v !== T;
			});
			for (var [n, i] of r) i.v !== T && !(n in e) && t.push(n);
			return t;
		},
		setPrototypeOf() {
			_e();
		}
	});
}
function Gt(e) {
	try {
		if (typeof e == "object" && e && ae in e) return e[ae];
	} catch {}
	return e;
}
function Kt(e, t) {
	return Object.is(Gt(e), Gt(t));
}
var qt, Jt, Yt, Xt;
function Zt() {
	if (qt === void 0) {
		qt = window, document, Jt = /Firefox/.test(navigator.userAgent);
		var e = Element.prototype, t = Node.prototype, n = Text.prototype;
		Yt = a(t, "firstChild").get, Xt = a(t, "nextSibling").get, u(e) && (e.__click = void 0, e.__className = void 0, e.__attributes = null, e.__style = void 0, e.__e = void 0), u(n) && (n.__t = void 0);
	}
}
function I(e = "") {
	return document.createTextNode(e);
}
/* @__NO_SIDE_EFFECTS__ */
function Qt(e) {
	return Yt.call(e);
}
/* @__NO_SIDE_EFFECTS__ */
function $t(e) {
	return Xt.call(e);
}
function L(e, t) {
	if (!E) return /* @__PURE__ */ Qt(e);
	var n = /* @__PURE__ */ Qt(D);
	if (n === null) n = D.appendChild(I());
	else if (t && n.nodeType !== 3) {
		var r = I();
		return n?.before(r), Ee(r), r;
	}
	return t && an(n), Ee(n), n;
}
function en(e, t = !1) {
	if (!E) {
		var n = /* @__PURE__ */ Qt(e);
		return n instanceof Comment && n.data === "" ? /* @__PURE__ */ $t(n) : n;
	}
	if (t) {
		if (D?.nodeType !== 3) {
			var r = I();
			return D?.before(r), Ee(r), r;
		}
		an(D);
	}
	return D;
}
function R(e, t = 1, n = !1) {
	let r = E ? D : e;
	for (var i; t--;) i = r, r = /* @__PURE__ */ $t(r);
	if (!E) return r;
	if (n) {
		if (r?.nodeType !== 3) {
			var a = I();
			return r === null ? i?.after(a) : r.before(a), Ee(a), a;
		}
		an(r);
	}
	return Ee(r), r;
}
function tn(e) {
	e.textContent = "";
}
function nn() {
	return !Pe || dt !== null ? !1 : (W.f & y) !== 0;
}
function rn(e, t, n) {
	let r = n ? { is: n } : void 0;
	return document.createElementNS(t ?? "http://www.w3.org/1999/xhtml", e, r);
}
function an(e) {
	if (e.nodeValue.length < 65536) return;
	let t = e.nextSibling;
	for (; t !== null && t.nodeType === 3;) t.remove(), e.nodeValue += t.nodeValue, t = e.nextSibling;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/misc.js
var on = !1;
function sn() {
	on || (on = !0, document.addEventListener("reset", (e) => {
		Promise.resolve().then(() => {
			if (!e.defaultPrevented) for (let t of e.target.elements) t.__on_r?.();
		});
	}, { capture: !0 }));
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/shared.js
function cn(e) {
	var t = H, n = W;
	U(null), Ln(null);
	try {
		return e();
	} finally {
		U(t), Ln(n);
	}
}
function ln(e, t, n, r = n) {
	e.addEventListener(t, () => cn(n));
	let i = e.__on_r;
	i ? e.__on_r = () => {
		i(), r(!0);
	} : e.__on_r = () => r(!0), sn();
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/effects.js
function un(e) {
	W === null && (H === null && pe(e), fe()), Pn && de(e);
}
function dn(e, t) {
	var n = t.last;
	n === null ? t.last = t.first = e : (n.next = e, e.prev = n, t.last = e);
}
function fn(e, t) {
	var n = W;
	n !== null && n.f & 8192 && (e |= _);
	var r = {
		ctx: k,
		deps: null,
		nodes: null,
		f: e | h | 512,
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
	if (e & 4) it === null ? ct.ensure().schedule(r) : it.push(r);
	else if (t !== null) {
		try {
			Zn(r);
		} catch (e) {
			throw V(r), e;
		}
		i.deps === null && i.teardown === null && i.nodes === null && i.first === i.last && !(i.f & 524288) && (i = i.first, e & 16 && e & 65536 && i !== null && (i.f |= x));
	}
	if (i !== null && (i.parent = n, n !== null && dn(i, n), H !== null && H.f & 2 && !(e & 64))) {
		var a = H;
		(a.effects ??= []).push(i);
	}
	return r;
}
function pn() {
	return H !== null && !In;
}
function mn(e) {
	let t = fn(8, null);
	return A(t, m), t.teardown = e, t;
}
function hn(e) {
	un("$effect");
	var t = W.f;
	if (!H && t & 32 && !(t & 32768)) {
		var n = k;
		(n.e ??= []).push(e);
	} else return gn(e);
}
function gn(e) {
	return fn(4 | ee, e);
}
function _n(e) {
	ct.ensure();
	let t = fn(64 | S, e);
	return (e = {}) => new Promise((n) => {
		e.outro ? Dn(t, () => {
			V(t), n(void 0);
		}) : (V(t), n(void 0));
	});
}
function vn(e) {
	return fn(4, e);
}
function yn(e) {
	return fn(re | S, e);
}
function bn(e, t = 0) {
	return fn(8 | t, e);
}
function z(e, t = [], n = [], r = []) {
	St(r, t, n, (t) => {
		fn(8, () => e(...t.map(J)));
	});
}
function xn(e, t = 0) {
	return fn(16 | t, e);
}
function B(e) {
	return fn(32 | S, e);
}
function Sn(e) {
	var t = e.teardown;
	if (t !== null) {
		let e = Pn, n = H;
		Fn(!0), U(null);
		try {
			t.call(null);
		} finally {
			Fn(e), U(n);
		}
	}
}
function Cn(e, t = !1) {
	var n = e.first;
	for (e.first = e.last = null; n !== null;) {
		let e = n.ac;
		e !== null && cn(() => {
			e.abort(w);
		});
		var r = n.next;
		n.f & 64 ? n.parent = null : V(n, t), n = r;
	}
}
function wn(e) {
	for (var t = e.first; t !== null;) {
		var n = t.next;
		t.f & 32 || V(t), t = n;
	}
}
function V(e, t = !0) {
	var n = !1;
	(t || e.f & 262144) && e.nodes !== null && e.nodes.end !== null && (Tn(e.nodes.start, e.nodes.end), n = !0), A(e, b), Cn(e, t && !n), Xn(e, 0);
	var r = e.nodes && e.nodes.t;
	if (r !== null) for (let e of r) e.stop();
	Sn(e), e.f ^= b, e.f |= v;
	var i = e.parent;
	i !== null && i.first !== null && En(e), e.next = e.prev = e.teardown = e.ctx = e.deps = e.fn = e.nodes = e.ac = null;
}
function Tn(e, t) {
	for (; e !== null;) {
		var n = e === t ? null : /* @__PURE__ */ $t(e);
		e.remove(), e = n;
	}
}
function En(e) {
	var t = e.parent, n = e.prev, r = e.next;
	n !== null && (n.next = r), r !== null && (r.prev = n), t !== null && (t.first === e && (t.first = r), t.last === e && (t.last = n));
}
function Dn(e, t, n = !0) {
	var r = [];
	On(e, r, !0);
	var i = () => {
		n && V(e), t && t();
	}, a = r.length;
	if (a > 0) {
		var o = () => --a || i();
		for (var s of r) s.out(o);
	} else i();
}
function On(e, t, n) {
	if (!(e.f & 8192)) {
		e.f ^= _;
		var r = e.nodes && e.nodes.t;
		if (r !== null) for (let e of r) (e.is_global || n) && t.push(e);
		for (var i = e.first; i !== null;) {
			var a = i.next, o = (i.f & 65536) != 0 || (i.f & 32) != 0 && (e.f & 16) != 0;
			On(i, t, o ? n : !1), i = a;
		}
	}
}
function kn(e) {
	An(e, !0);
}
function An(e, t) {
	if (e.f & 8192) {
		e.f ^= _, e.f & 1024 || (A(e, h), ct.ensure().schedule(e));
		for (var n = e.first; n !== null;) {
			var r = n.next, i = (n.f & 65536) != 0 || (n.f & 32) != 0;
			An(n, i ? t : !1), n = r;
		}
		var a = e.nodes && e.nodes.t;
		if (a !== null) for (let e of a) (e.is_global || t) && e.in();
	}
}
function jn(e, t) {
	if (e.nodes) for (var n = e.nodes.start, r = e.nodes.end; n !== null;) {
		var i = n === r ? null : /* @__PURE__ */ $t(n);
		t.append(n), n = i;
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/legacy.js
var Mn = null, Nn = !1, Pn = !1;
function Fn(e) {
	Pn = e;
}
var H = null, In = !1;
function U(e) {
	H = e;
}
var W = null;
function Ln(e) {
	W = e;
}
var G = null;
function Rn(e) {
	H !== null && (!Pe || H.f & 2) && (G === null ? G = [e] : G.push(e));
}
var K = null, q = 0, zn = null;
function Bn(e) {
	zn = e;
}
var Vn = 1, Hn = 0, Un = Hn;
function Wn(e) {
	Un = e;
}
function Gn() {
	return ++Vn;
}
function Kn(e) {
	var t = e.f;
	if (t & 2048) return !0;
	if (t & 2 && (e.f &= ~C), t & 4096) {
		for (var n = e.deps, r = n.length, i = 0; i < r; i++) {
			var a = n[i];
			if (Kn(a) && Mt(a), a.wv > e.wv) return !0;
		}
		t & 512 && M === null && A(e, m);
	}
	return !1;
}
function qn(e, t, r = !0) {
	var i = e.reactions;
	if (i !== null && !(!Pe && G !== null && n.call(G, e))) for (var a = 0; a < i.length; a++) {
		var o = i[a];
		o.f & 2 ? qn(o, t, !1) : t === o && (r ? A(o, h) : o.f & 1024 && A(o, g), ht(o));
	}
}
function Jn(e) {
	var t = K, n = q, r = zn, i = H, a = G, o = k, s = In, c = Un, l = e.f;
	K = null, q = 0, zn = null, H = l & 96 ? null : e, G = null, Ie(e.ctx), In = !1, Un = ++Hn, e.ac !== null && (cn(() => {
		e.ac.abort(w);
	}), e.ac = null);
	try {
		e.f |= ne;
		var u = e.fn, d = u();
		e.f |= y;
		var f = e.deps, p = j?.is_fork;
		if (K !== null) {
			var m;
			if (p || Xn(e, q), f !== null && q > 0) for (f.length = q + K.length, m = 0; m < K.length; m++) f[q + m] = K[m];
			else e.deps = f = K;
			if (pn() && e.f & 512) for (m = q; m < f.length; m++) (f[m].reactions ??= []).push(e);
		} else !p && f !== null && q < f.length && (Xn(e, q), f.length = q);
		if (ze() && zn !== null && !In && f !== null && !(e.f & 6146)) for (m = 0; m < zn.length; m++) qn(zn[m], e);
		if (i !== null && i !== e) {
			if (Hn++, i.deps !== null) for (let e = 0; e < n; e += 1) i.deps[e].rv = Hn;
			if (t !== null) for (let e of t) e.rv = Hn;
			zn !== null && (r === null ? r = zn : r.push(...zn));
		}
		return e.f & 8388608 && (e.f ^= ie), d;
	} catch (e) {
		return We(e);
	} finally {
		e.f ^= ne, K = t, q = n, zn = r, H = i, G = a, Ie(o), In = s, Un = c;
	}
}
function Yn(e, r) {
	let i = r.reactions;
	if (i !== null) {
		var a = t.call(i, e);
		if (a !== -1) {
			var o = i.length - 1;
			o === 0 ? i = r.reactions = null : (i[a] = i[o], i.pop());
		}
	}
	if (i === null && r.f & 2 && (K === null || !n.call(K, r))) {
		var s = r;
		s.f & 512 && (s.f ^= 512, s.f &= ~C), qe(s), Nt(s), Xn(s, 0);
	}
}
function Xn(e, t) {
	var n = e.deps;
	if (n !== null) for (var r = t; r < n.length; r++) Yn(e, n[r]);
}
function Zn(e) {
	var t = e.f;
	if (!(t & 16384)) {
		A(e, m);
		var n = W, r = Nn;
		W = e, Nn = !0;
		try {
			t & 16777232 ? wn(e) : Cn(e), Sn(e);
			var i = Jn(e);
			e.teardown = typeof i == "function" ? i : null, e.wv = Vn;
		} finally {
			Nn = r, W = n;
		}
	}
}
async function Qn() {
	if (Pe) return new Promise((e) => {
		requestAnimationFrame(() => e()), setTimeout(() => e());
	});
	await Promise.resolve(), lt();
}
function J(e) {
	var t = (e.f & 2) != 0;
	if (Mn?.add(e), H !== null && !In && !(W !== null && W.f & 16384) && (G === null || !n.call(G, e))) {
		var r = H.deps;
		if (H.f & 2097152) e.rv < Hn && (e.rv = Hn, K === null && r !== null && r[q] === e ? q++ : K === null ? K = [e] : K.push(e));
		else {
			(H.deps ??= []).push(e);
			var i = e.reactions;
			i === null ? e.reactions = [H] : n.call(i, H) || i.push(H);
		}
	}
	if (Pn && It.has(e)) return It.get(e);
	if (t) {
		var a = e;
		if (Pn) {
			var o = a.v;
			return (!(a.f & 1024) && a.reactions !== null || er(a)) && (o = jt(a)), It.set(a, o), o;
		}
		var s = (a.f & 512) == 0 && !In && H !== null && (Nn || (H.f & 512) != 0), c = (a.f & y) === 0;
		Kn(a) && (s && (a.f |= 512), Mt(a)), s && !c && (Pt(a), $n(a));
	}
	if (M?.has(e)) return M.get(e);
	if (e.f & 8388608) throw e.v;
	return e.v;
}
function $n(e) {
	if (e.f |= 512, e.deps !== null) for (let t of e.deps) (t.reactions ??= []).push(e), t.f & 2 && !(t.f & 512) && (Pt(t), $n(t));
}
function er(e) {
	if (e.v === T) return !0;
	if (e.deps === null) return !1;
	for (let t of e.deps) if (It.has(t) || t.f & 2 && er(t)) return !0;
	return !1;
}
function tr(e) {
	var t = In;
	try {
		return In = !0, e();
	} finally {
		In = t;
	}
}
[.../* @__PURE__ */ "allowfullscreen.async.autofocus.autoplay.checked.controls.default.disabled.formnovalidate.indeterminate.inert.ismap.loop.multiple.muted.nomodule.novalidate.open.playsinline.readonly.required.reversed.seamless.selected.webkitdirectory.defer.disablepictureinpicture.disableremoteplayback".split(".")];
var nr = ["touchstart", "touchmove"];
function rr(e) {
	return nr.includes(e);
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/events.js
var ir = Symbol("events"), ar = /* @__PURE__ */ new Set(), or = /* @__PURE__ */ new Set();
function sr(e, t, n, r = {}) {
	function i(e) {
		if (r.capture || dr.call(t, e), !e.cancelBubble) return cn(() => n?.call(this, e));
	}
	return e.startsWith("pointer") || e.startsWith("touch") || e === "wheel" ? He(() => {
		t.addEventListener(e, i, r);
	}) : t.addEventListener(e, i, r), i;
}
function cr(e, t, n, r, i) {
	var a = {
		capture: r,
		passive: i
	}, o = sr(e, t, n, a);
	(t === document.body || t === window || t === document || t instanceof HTMLMediaElement) && mn(() => {
		t.removeEventListener(e, o, a);
	});
}
function Y(e, t, n) {
	(t[ir] ??= {})[e] = n;
}
function lr(e) {
	for (var t = 0; t < e.length; t++) ar.add(e[t]);
	for (var n of or) n(e);
}
var ur = null;
function dr(e) {
	var t = this, n = t.ownerDocument, r = e.type, a = e.composedPath?.() || [], o = a[0] || e.target;
	ur = e;
	var s = 0, c = ur === e && e[ir];
	if (c) {
		var l = a.indexOf(c);
		if (l !== -1 && (t === document || t === window)) {
			e[ir] = t;
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
		var d = H, f = W;
		U(null), Ln(null);
		try {
			for (var p, m = []; o !== null;) {
				var h = o.assignedSlot || o.parentNode || o.host || null;
				try {
					var g = o[ir]?.[r];
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
			e[ir] = t, delete e.currentTarget, U(d), Ln(f);
		}
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/reconciler.js
var fr = globalThis?.window?.trustedTypes && /* @__PURE__ */ globalThis.window.trustedTypes.createPolicy("svelte-trusted-html", { createHTML: (e) => e });
function pr(e) {
	return fr?.createHTML(e) ?? e;
}
function mr(e) {
	var t = rn("template");
	return t.innerHTML = pr(e.replaceAll("<!>", "<!---->")), t.content;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/template.js
function hr(e, t) {
	var n = W;
	n.nodes === null && (n.nodes = {
		start: e,
		end: t,
		a: null,
		t: null
	});
}
/* @__NO_SIDE_EFFECTS__ */
function X(e, t) {
	var n = (t & 1) != 0, r = (t & 2) != 0, i, a = !e.startsWith("<!>");
	return () => {
		if (E) return hr(D, null), D;
		i === void 0 && (i = mr(a ? e : "<!>" + e), n || (i = /* @__PURE__ */ Qt(i)));
		var t = r || Jt ? document.importNode(i, !0) : i.cloneNode(!0);
		if (n) {
			var o = /* @__PURE__ */ Qt(t), s = t.lastChild;
			hr(o, s);
		} else hr(t, t);
		return t;
	};
}
function gr() {
	if (E) return hr(D, null), D;
	var e = document.createDocumentFragment(), t = document.createComment(""), n = I();
	return e.append(t, n), hr(t, n), e;
}
function Z(e, t) {
	if (E) {
		var n = W;
		(!(n.f & 32768) || n.nodes.end === null) && (n.nodes.end = D), De();
		return;
	}
	e !== null && e.before(t);
}
function Q(e, t) {
	var n = t == null ? "" : typeof t == "object" ? `${t}` : t;
	n !== (e.__t ??= e.nodeValue) && (e.__t = n, e.nodeValue = `${n}`);
}
function _r(e, t) {
	return yr(e, t);
}
var vr = /* @__PURE__ */ new Map();
function yr(e, { target: t, anchor: n, props: i = {}, events: a, context: o, intro: s = !0, transformError: c }) {
	Zt();
	var l = void 0, u = _n(() => {
		var s = n ?? t.appendChild(I());
		bt(s, { pending: () => {} }, (t) => {
			Le({});
			var n = k;
			if (o && (n.c = o), a && (i.$$events = a), E && hr(t, null), l = e(t, i) || {}, E && (W.nodes.end = D, D === null || D.nodeType !== 8 || D.data !== "]")) throw Se(), be;
			Re();
		}, c);
		var u = /* @__PURE__ */ new Set(), d = (e) => {
			for (var n = 0; n < e.length; n++) {
				var r = e[n];
				if (!u.has(r)) {
					u.add(r);
					var i = rr(r);
					for (let e of [t, document]) {
						var a = vr.get(e);
						a === void 0 && (a = /* @__PURE__ */ new Map(), vr.set(e, a));
						var o = a.get(r);
						o === void 0 ? (e.addEventListener(r, dr, { passive: i }), a.set(r, 1)) : a.set(r, o + 1);
					}
				}
			}
		};
		return d(r(ar)), or.add(d), () => {
			for (var e of u) for (let n of [t, document]) {
				var r = vr.get(n), i = r.get(e);
				--i == 0 ? (n.removeEventListener(e, dr), r.delete(e), r.size === 0 && vr.delete(n)) : r.set(e, i);
			}
			or.delete(d), s !== n && s.parentNode?.removeChild(s);
		};
	});
	return br.set(l, u), l;
}
var br = /* @__PURE__ */ new WeakMap(), xr = class {
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
			if (n) kn(n), this.#r.delete(t);
			else {
				var r = this.#n.get(t);
				r && (this.#t.set(t, r.effect), this.#n.delete(t), r.fragment.lastChild.remove(), this.anchor.before(r.fragment), n = r.effect);
			}
			for (let [t, n] of this.#e) {
				if (this.#e.delete(t), t === e) break;
				let r = this.#n.get(n);
				r && (V(r.effect), this.#n.delete(n));
			}
			for (let [e, r] of this.#t) {
				if (e === t || this.#r.has(e)) continue;
				let i = () => {
					if (Array.from(this.#e.values()).includes(e)) {
						var t = document.createDocumentFragment();
						jn(r, t), t.append(I()), this.#n.set(e, {
							effect: r,
							fragment: t
						});
					} else V(r);
					this.#r.delete(e), this.#t.delete(e);
				};
				this.#i || !n ? (this.#r.add(e), Dn(r, i, !1)) : i();
			}
		}
	};
	#o = (e) => {
		this.#e.delete(e);
		let t = Array.from(this.#e.values());
		for (let [e, n] of this.#n) t.includes(e) || (V(n.effect), this.#n.delete(e));
	};
	ensure(e, t) {
		var n = j, r = nn();
		if (t && !this.#t.has(e) && !this.#n.has(e)) if (r) {
			var i = document.createDocumentFragment(), a = I();
			i.append(a), this.#n.set(e, {
				effect: B(() => t(a)),
				fragment: i
			});
		} else this.#t.set(e, B(() => t(this.anchor)));
		if (this.#e.set(n, e), r) {
			for (let [t, r] of this.#t) t === e ? n.unskip_effect(r) : n.skip_effect(r);
			for (let [t, r] of this.#n) t === e ? n.unskip_effect(r.effect) : n.skip_effect(r.effect);
			n.oncommit(this.#a), n.ondiscard(this.#o);
		} else E && (this.anchor = D), this.#a(n);
	}
};
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/if.js
function $(e, t, n = !1) {
	var r;
	E && (r = D, De());
	var i = new xr(e), a = n ? x : 0;
	function o(e, t) {
		if (E) {
			var n = Ae(r);
			if (e !== parseInt(n.substring(1))) {
				var a = ke();
				Ee(a), i.anchor = a, Te(!1), i.ensure(e, t), Te(!0);
				return;
			}
		}
		i.ensure(e, t);
	}
	xn(() => {
		var e = !1;
		t((t, n = 0) => {
			e = !0, o(n, t);
		}), e || o(-1, null);
	}, a);
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/each.js
function Sr(e, t) {
	return t;
}
function Cr(e, t, n) {
	for (var i = [], a = t.length, o, s = t.length, c = 0; c < a; c++) {
		let n = t[c];
		Dn(n, () => {
			if (o) {
				if (o.pending.delete(n), o.done.add(n), o.pending.size === 0) {
					var t = e.outrogroups;
					wr(e, r(o.done)), t.delete(o), t.size === 0 && (e.outrogroups = null);
				}
			} else --s;
		}, !1);
	}
	if (s === 0) {
		var l = i.length === 0 && n !== null;
		if (l) {
			var u = n, d = u.parentNode;
			tn(d), d.append(u), e.items.clear();
		}
		wr(e, t, !l);
	} else o = {
		pending: new Set(t),
		done: /* @__PURE__ */ new Set()
	}, (e.outrogroups ??= /* @__PURE__ */ new Set()).add(o);
}
function wr(e, t, n = !0) {
	var r;
	if (e.pending.size > 0) {
		r = /* @__PURE__ */ new Set();
		for (let t of e.pending.values()) for (let n of t) r.add(e.items.get(n).e);
	}
	for (var i = 0; i < t.length; i++) {
		var a = t[i];
		r?.has(a) ? (a.f |= te, jn(a, document.createDocumentFragment())) : V(t[i], n);
	}
}
var Tr;
function Er(t, n, i, a, o, s = null) {
	var c = t, l = /* @__PURE__ */ new Map();
	if (n & 4) {
		var u = t;
		c = E ? Ee(/* @__PURE__ */ Qt(u)) : u.appendChild(I());
	}
	E && De();
	var d = null, f = /* @__PURE__ */ Ot(() => {
		var t = i();
		return e(t) ? t : t == null ? [] : r(t);
	}), p, m = /* @__PURE__ */ new Map(), h = !0;
	function g(e) {
		v.effect.f & 16384 || (v.pending.delete(e), v.fallback = d, Or(v, p, c, n, a), d !== null && (p.length === 0 ? d.f & 33554432 ? (d.f ^= te, Ar(d, null, c)) : kn(d) : Dn(d, () => {
			d = null;
		})));
	}
	function _(e) {
		v.pending.delete(e);
	}
	var v = {
		effect: xn(() => {
			p = J(f);
			var e = p.length;
			let t = !1;
			E && Ae(c) === "[!" != (e === 0) && (c = ke(), Ee(c), Te(!1), t = !0);
			for (var r = /* @__PURE__ */ new Set(), u = j, v = nn(), y = 0; y < e; y += 1) {
				E && D.nodeType === 8 && D.data === "]" && (c = D, t = !0, Te(!1));
				var b = p[y], x = a(b, y), S = h ? null : l.get(x);
				S ? (S.v && Bt(S.v, b), S.i && Bt(S.i, y), v && u.unskip_effect(S.e)) : (S = kr(l, h ? c : Tr ??= I(), b, x, y, o, n, i), h || (S.e.f |= te), l.set(x, S)), r.add(x);
			}
			if (e === 0 && s && !d && (h ? d = B(() => s(c)) : (d = B(() => s(Tr ??= I())), d.f |= te)), e > r.size && ue("", "", ""), E && e > 0 && Ee(ke()), !h) if (m.set(u, r), v) {
				for (let [e, t] of l) r.has(e) || u.skip_effect(t.e);
				u.oncommit(g), u.ondiscard(_);
			} else g(u);
			t && Te(!0), J(f);
		}),
		flags: n,
		items: l,
		pending: m,
		outrogroups: null,
		fallback: d
	};
	h = !1, E && (c = D);
}
function Dr(e) {
	for (; e !== null && !(e.f & 32);) e = e.next;
	return e;
}
function Or(e, t, n, i, a) {
	var o = (i & 8) != 0, s = t.length, c = e.items, l = Dr(e.effect.first), u, d = null, f, p = [], m = [], h, g, _, v;
	if (o) for (v = 0; v < s; v += 1) h = t[v], g = a(h, v), _ = c.get(g).e, _.f & 33554432 || (_.nodes?.a?.measure(), (f ??= /* @__PURE__ */ new Set()).add(_));
	for (v = 0; v < s; v += 1) {
		if (h = t[v], g = a(h, v), _ = c.get(g).e, e.outrogroups !== null) for (let t of e.outrogroups) t.pending.delete(_), t.done.delete(_);
		if (_.f & 33554432) if (_.f ^= te, _ === l) Ar(_, null, n);
		else {
			var y = d ? d.next : l;
			_ === e.effect.last && (e.effect.last = _.prev), _.prev && (_.prev.next = _.next), _.next && (_.next.prev = _.prev), jr(e, d, _), jr(e, _, y), Ar(_, y, n), d = _, p = [], m = [], l = Dr(d.next);
			continue;
		}
		if (_.f & 8192 && (kn(_), o && (_.nodes?.a?.unfix(), (f ??= /* @__PURE__ */ new Set()).delete(_))), _ !== l) {
			if (u !== void 0 && u.has(_)) {
				if (p.length < m.length) {
					var b = m[0], x;
					d = b.prev;
					var S = p[0], ee = p[p.length - 1];
					for (x = 0; x < p.length; x += 1) Ar(p[x], b, n);
					for (x = 0; x < m.length; x += 1) u.delete(m[x]);
					jr(e, S.prev, ee.next), jr(e, d, S), jr(e, ee, b), l = b, d = ee, --v, p = [], m = [];
				} else u.delete(_), Ar(_, l, n), jr(e, _.prev, _.next), jr(e, _, d === null ? e.effect.first : d.next), jr(e, d, _), d = _;
				continue;
			}
			for (p = [], m = []; l !== null && l !== _;) (u ??= /* @__PURE__ */ new Set()).add(l), m.push(l), l = Dr(l.next);
			if (l === null) continue;
		}
		_.f & 33554432 || p.push(_), d = _, l = Dr(_.next);
	}
	if (e.outrogroups !== null) {
		for (let t of e.outrogroups) t.pending.size === 0 && (wr(e, r(t.done)), e.outrogroups?.delete(t));
		e.outrogroups.size === 0 && (e.outrogroups = null);
	}
	if (l !== null || u !== void 0) {
		var C = [];
		if (u !== void 0) for (_ of u) _.f & 8192 || C.push(_);
		for (; l !== null;) !(l.f & 8192) && l !== e.fallback && C.push(l), l = Dr(l.next);
		var ne = C.length;
		if (ne > 0) {
			var re = i & 4 && s === 0 ? n : null;
			if (o) {
				for (v = 0; v < ne; v += 1) C[v].nodes?.a?.measure();
				for (v = 0; v < ne; v += 1) C[v].nodes?.a?.fix();
			}
			Cr(e, C, re);
		}
	}
	o && He(() => {
		if (f !== void 0) for (_ of f) _.nodes?.a?.apply();
	});
}
function kr(e, t, n, r, i, a, o, s) {
	var c = o & 1 ? o & 16 ? Rt(n) : /* @__PURE__ */ zt(n, !1, !1) : null, l = o & 2 ? Rt(i) : null;
	return {
		v: c,
		i: l,
		e: B(() => (a(t, c ?? n, l ?? i, s), () => {
			e.delete(r);
		}))
	};
}
function Ar(e, t, n) {
	if (e.nodes) for (var r = e.nodes.start, i = e.nodes.end, a = t && !(t.f & 33554432) ? t.nodes.start : n; r !== null;) {
		var o = /* @__PURE__ */ $t(r);
		if (a.before(r), r === i) return;
		r = o;
	}
}
function jr(e, t, n) {
	t === null ? e.effect.first = n : t.next = n, n === null ? e.effect.last = t : n.prev = t;
}
//#endregion
//#region node_modules/clsx/dist/clsx.mjs
function Mr(e) {
	var t, n, r = "";
	if (typeof e == "string" || typeof e == "number") r += e;
	else if (typeof e == "object") if (Array.isArray(e)) {
		var i = e.length;
		for (t = 0; t < i; t++) e[t] && (n = Mr(e[t])) && (r && (r += " "), r += n);
	} else for (n in e) e[n] && (r && (r += " "), r += n);
	return r;
}
function Nr() {
	for (var e, t, n = 0, r = "", i = arguments.length; n < i; n++) (e = arguments[n]) && (t = Mr(e)) && (r && (r += " "), r += t);
	return r;
}
//#endregion
//#region node_modules/svelte/src/internal/shared/attributes.js
function Pr(e) {
	return typeof e == "object" ? Nr(e) : e ?? "";
}
var Fr = [..." 	\n\r\f\xA0\v﻿"];
function Ir(e, t, n) {
	var r = e == null ? "" : "" + e;
	if (t && (r = r ? r + " " + t : t), n) {
		for (var i of Object.keys(n)) if (n[i]) r = r ? r + " " + i : i;
		else if (r.length) for (var a = i.length, o = 0; (o = r.indexOf(i, o)) >= 0;) {
			var s = o + a;
			(o === 0 || Fr.includes(r[o - 1])) && (s === r.length || Fr.includes(r[s])) ? r = (o === 0 ? "" : r.substring(0, o)) + r.substring(s + 1) : o = s;
		}
	}
	return r === "" ? null : r;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/class.js
function Lr(e, t, n, r, i, a) {
	var o = e.__className;
	if (E || o !== n || o === void 0) {
		var s = Ir(n, r, a);
		(!E || s !== e.getAttribute("class")) && (s == null ? e.removeAttribute("class") : t ? e.className = s : e.setAttribute("class", s)), e.__className = n;
	} else if (a && i !== a) for (var c in a) {
		var l = !!a[c];
		(i == null || l !== !!i[c]) && e.classList.toggle(c, l);
	}
	return a;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/select.js
function Rr(t, n, r = !1) {
	if (t.multiple) {
		if (n == null) return;
		if (!e(n)) return Ce();
		for (var i of t.options) i.selected = n.includes(Br(i));
		return;
	}
	for (i of t.options) if (Kt(Br(i), n)) {
		i.selected = !0;
		return;
	}
	(!r || n !== void 0) && (t.selectedIndex = -1);
}
function zr(e) {
	var t = new MutationObserver(() => {
		Rr(e, e.__value);
	});
	t.observe(e, {
		childList: !0,
		subtree: !0,
		attributes: !0,
		attributeFilter: ["value"]
	}), mn(() => {
		t.disconnect();
	});
}
function Br(e) {
	return "__value" in e ? e.__value : e.value;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/attributes.js
var Vr = Symbol("is custom element"), Hr = Symbol("is html"), Ur = ce ? "link" : "LINK", Wr = ce ? "progress" : "PROGRESS";
function Gr(e) {
	if (E) {
		var t = !1, n = () => {
			if (!t) {
				if (t = !0, e.hasAttribute("value")) {
					var n = e.value;
					qr(e, "value", null), e.value = n;
				}
				if (e.hasAttribute("checked")) {
					var r = e.checked;
					qr(e, "checked", null), e.checked = r;
				}
			}
		};
		e.__on_r = n, He(n), sn();
	}
}
function Kr(e, t) {
	var n = Jr(e);
	n.value === (n.value = t ?? void 0) || e.value === t && (t !== 0 || e.nodeName !== Wr) || (e.value = t ?? "");
}
function qr(e, t, n, r) {
	var i = Jr(e);
	E && (i[t] = e.getAttribute(t), t === "src" || t === "srcset" || t === "href" && e.nodeName === Ur) || i[t] !== (i[t] = n) && (t === "loading" && (e[se] = n), n == null ? e.removeAttribute(t) : typeof n != "string" && Xr(e).includes(t) ? e[t] = n : e.setAttribute(t, n));
}
function Jr(e) {
	return e.__attributes ??= {
		[Vr]: e.nodeName.includes("-"),
		[Hr]: e.namespaceURI === xe
	};
}
var Yr = /* @__PURE__ */ new Map();
function Xr(e) {
	var t = e.getAttribute("is") || e.nodeName, n = Yr.get(t);
	if (n) return n;
	Yr.set(t, n = []);
	for (var r, i = e, a = Element.prototype; a !== i;) {
		for (var s in r = o(i), r) r[s].set && n.push(s);
		i = l(i);
	}
	return n;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/input.js
function Zr(e, t, n = t) {
	var r = /* @__PURE__ */ new WeakSet();
	ln(e, "input", async (i) => {
		var a = i ? e.defaultValue : e.value;
		if (a = Qr(e) ? $r(a) : a, n(a), j !== null && r.add(j), await Qn(), a !== (a = t())) {
			var o = e.selectionStart, s = e.selectionEnd, c = e.value.length;
			if (e.value = a ?? "", s !== null) {
				var l = e.value.length;
				o === s && s === c && l > c ? (e.selectionStart = l, e.selectionEnd = l) : (e.selectionStart = o, e.selectionEnd = Math.min(s, l));
			}
		}
	}), (E && e.defaultValue !== e.value || tr(t) == null && e.value) && (n(Qr(e) ? $r(e.value) : e.value), j !== null && r.add(j)), bn(() => {
		var n = t();
		if (e === document.activeElement) {
			var i = Pe ? et : j;
			if (r.has(i)) return;
		}
		Qr(e) && n === $r(e.value) || e.type === "date" && !n && !e.value || n !== e.value && (e.value = n ?? "");
	});
}
function Qr(e) {
	var t = e.type;
	return t === "number" || t === "range";
}
function $r(e) {
	return e === "" ? null : +e;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/this.js
function ei(e, t) {
	return e === t || e?.[ae] === t;
}
function ti(e = {}, t, n, r) {
	var i = k.r, a = W;
	return vn(() => {
		var o, s;
		return bn(() => {
			o = s, s = r?.() || [], tr(() => {
				e !== n(...s) && (t(e, ...s), o && ei(n(...o), e) && t(null, ...o));
			});
		}), () => {
			let r = a;
			for (; r !== i && r.parent !== null && r.parent.f & 33554432;) r = r.parent;
			let o = () => {
				s && ei(n(...s), e) && t(null, ...s);
			}, c = r.teardown;
			r.teardown = () => {
				o(), c?.();
			};
		};
	}), e;
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/props.js
function ni(e, t, n, r) {
	var i = !Fe || (n & 2) != 0, o = (n & 8) != 0, s = (n & 16) != 0, c = r, l = !0, u = () => (l && (l = !1, c = s ? tr(r) : r), c);
	let d;
	if (o) {
		var f = ae in e || oe in e;
		d = a(e, t)?.set ?? (f && t in e ? (n) => e[t] = n : void 0);
	}
	var p, m = !1;
	o ? [p, m] = Qe(() => e[t]) : p = e[t], p === void 0 && r !== void 0 && (p = u(), d && (i && he(t), d(p)));
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
	var _ = !1, v = (n & 1 ? Et : Ot)(() => (_ = !1, h()));
	o && J(v);
	var y = W;
	return (function(e, t) {
		if (arguments.length > 0) {
			let n = t ? J(v) : i && o ? Wt(e) : e;
			return F(v, n), _ = !0, c !== void 0 && (c = n), e;
		}
		return Pn && _ || y.f & 16384 ? v.v : J(v);
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
var di = /* @__PURE__ */ X("<option> </option>"), fi = /* @__PURE__ */ X("<span class=\"block text-xs text-slate-400 mt-0.5 px-1\"> </span>"), pi = /* @__PURE__ */ X("<tr class=\"border-b border-slate-100 hover:bg-slate-50 text-sm group\"><td class=\"px-1 py-1 w-24\"><select></select></td><td class=\"px-1 py-1\"><input type=\"text\" class=\"w-full px-1 py-0.5 text-slate-800 bg-transparent border-0 rounded\n				focus:ring-2 focus:ring-blue-400 focus:bg-white\" placeholder=\"Item name\"/> <!></td><td class=\"px-1 py-1 w-20\"><input type=\"number\" step=\"any\" min=\"0\" class=\"w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 rounded\n				focus:ring-2 focus:ring-blue-400 focus:bg-white\"/></td><td class=\"px-1 py-1 w-16\"><input type=\"text\" class=\"w-full text-center text-slate-500 px-1 py-0.5 bg-transparent border-0 rounded\n				focus:ring-2 focus:ring-blue-400 focus:bg-white text-sm\" placeholder=\"ea\"/></td><td class=\"px-1 py-1 w-24\"><input type=\"number\" step=\"0.01\" min=\"0\" class=\"w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 rounded\n				focus:ring-2 focus:ring-blue-400 focus:bg-white\"/></td><td class=\"px-2 py-1.5 text-right font-mono text-slate-400 w-16 text-xs\"> </td><td class=\"px-2 py-1.5 text-right font-mono text-slate-500 w-24 text-xs\"> </td><td class=\"px-2 py-1.5 text-right font-mono font-medium w-24\"> </td><td class=\"px-1 py-1 w-8\"><button class=\"opacity-0 group-hover:opacity-100 text-slate-400 hover:text-red-500 transition-opacity p-0.5\" title=\"Delete item\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></button></td></tr>");
function mi(e, t) {
	Le(t, !0);
	let n = ni(t, "item", 7), r = [
		"materials",
		"labor",
		"equipment",
		"subs",
		"other"
	], i = {
		materials: "bg-blue-100 text-blue-700",
		labor: "bg-amber-100 text-amber-700",
		equipment: "bg-purple-100 text-purple-700",
		subs: "bg-green-100 text-green-700",
		other: "bg-slate-100 text-slate-600"
	}, a = /* @__PURE__ */ N(() => ri(n().category_type, t.globals, t.markupOverrides, t.markupEnabled)), o = /* @__PURE__ */ N(() => ii(n().unit_price, J(a))), s = /* @__PURE__ */ N(() => ai(n().quantity, n().unit_price, J(a))), c = /* @__PURE__ */ N(() => i[n().category_type] || i.other);
	function l() {
		t.onchange?.();
	}
	function u(e) {
		let t = parseFloat(e.target.value);
		isNaN(t) || (n().quantity = t, l());
	}
	function d(e) {
		let t = parseFloat(e.target.value);
		isNaN(t) || (n().unit_price = t, n().price_override = !0, l());
	}
	function f(e) {
		n().item_name = e.target.value, l();
	}
	function p(e) {
		n().unit = e.target.value, l();
	}
	function m(e) {
		n().category_type = e.target.value, l();
	}
	function h() {
		t.ondelete?.(n().id);
	}
	var g = pi(), _ = L(g), v = L(_);
	Er(v, 21, () => r, Sr, (e, t) => {
		var n = di(), r = L(n, !0);
		O(n);
		var i = {};
		z(() => {
			Q(r, J(t)), i !== (i = J(t)) && (n.value = (n.__value = J(t)) ?? "");
		}), Z(e, n);
	}), O(v);
	var y;
	zr(v), O(_);
	var b = R(_), x = L(b);
	Gr(x);
	var S = R(x, 2), ee = (e) => {
		var t = fi(), r = L(t, !0);
		O(t), z(() => Q(r, n().description)), Z(e, t);
	};
	$(S, (e) => {
		n().description && e(ee);
	}), O(b);
	var te = R(b), C = L(te);
	Gr(C), O(te);
	var ne = R(te), re = L(ne);
	Gr(re), O(ne);
	var ie = R(ne), ae = L(ie);
	Gr(ae), O(ie);
	var oe = R(ie), se = L(oe);
	O(oe);
	var w = R(oe), ce = L(w, !0);
	O(w);
	var le = R(w), ue = L(le, !0);
	O(le);
	var de = R(le), fe = L(de);
	O(de), O(g), z((e, t) => {
		Lr(v, 1, `w-full text-xs px-1 py-1 rounded border-0 bg-transparent font-medium cursor-pointer
				focus:ring-2 focus:ring-blue-400 focus:bg-white ${J(c) ?? ""}`), y !== (y = n().category_type) && (v.value = (v.__value = n().category_type) ?? "", Rr(v, n().category_type)), Kr(x, n().item_name), Kr(C, n().quantity), Kr(re, n().unit), Kr(ae, n().unit_price), Q(se, `${J(a) ?? ""}%`), Q(ce, e), Q(ue, t);
	}, [() => li(J(o)), () => li(J(s))]), Y("change", v, m), Y("input", x, f), Y("input", C, u), Y("input", re, p), Y("input", ae, d), Y("click", fe, h), Z(e, g), Re();
}
lr([
	"change",
	"input",
	"click"
]);
//#endregion
//#region src/lib/Autocomplete.svelte
var hi = /* @__PURE__ */ X("<button><div><span class=\"text-slate-800\"> </span> <span class=\"text-xs text-slate-400 ml-2\"> </span></div> <div class=\"flex items-center gap-2 text-xs\"><span class=\"font-mono text-slate-500\"> </span> <span class=\"text-slate-400\"> </span></div></button>"), gi = /* @__PURE__ */ X("<div class=\"absolute top-full left-0 right-0 mt-1 bg-white border border-slate-200 rounded-lg shadow-lg z-50 max-h-64 overflow-y-auto\"></div>"), _i = /* @__PURE__ */ X("<div class=\"absolute top-full left-0 right-0 mt-1 bg-white border border-slate-200 rounded-lg shadow-lg z-50\"><button class=\"w-full text-left px-3 py-2 text-sm text-slate-500 hover:bg-slate-50\">Add custom item: <span class=\"font-medium text-slate-700\"> </span></button></div>"), vi = /* @__PURE__ */ X("<div class=\"relative\"><input type=\"text\" placeholder=\"Search items or type a name...\" class=\"w-full px-3 py-2 text-sm border border-slate-300 rounded-lg\n			focus:ring-2 focus:ring-blue-400 focus:border-blue-400 outline-none\"/> <!></div>");
function yi(e, t) {
	Le(t, !0);
	let n = ni(t, "materialsDb", 19, () => []), r = ni(t, "ratesDb", 19, () => []), i = ni(t, "categoryType", 3, "materials"), a = /* @__PURE__ */ P(""), o = /* @__PURE__ */ P(0), s, c = /* @__PURE__ */ N(() => {
		if (J(a).length < 1) return [];
		let e = J(a).toLowerCase(), t;
		return t = i() === "materials" ? n().map((e) => ({
			id: e.id,
			name: e.name,
			category: e.category,
			unit: e.default_unit,
			price: e.default_price,
			type: "materials",
			source: "material"
		})) : r().filter((e) => i() === "labor" ? e.type === "labor" : i() === "equipment" ? e.type === "equipment" : !0).map((e) => ({
			id: e.id,
			name: e.name,
			category: e.category,
			unit: e.default_unit,
			price: e.default_price,
			type: i(),
			source: "rate"
		})), t.filter((t) => t.name.toLowerCase().includes(e) || t.category.toLowerCase().includes(e)).slice(0, 15);
	});
	function l(e) {
		e.key === "ArrowDown" ? (e.preventDefault(), F(o, Math.min(J(o) + 1, J(c).length - 1), !0)) : e.key === "ArrowUp" ? (e.preventDefault(), F(o, Math.max(J(o) - 1, 0), !0)) : e.key === "Enter" ? (e.preventDefault(), J(c).length > 0 && J(o) < J(c).length ? u(J(c)[J(o)]) : J(a).trim() && d()) : e.key === "Escape" && (e.preventDefault(), t.oncancel?.());
	}
	function u(e) {
		t.onselect?.({
			item_name: e.name,
			unit: e.unit,
			unit_price: e.price,
			category_type: e.type,
			is_custom: !1,
			material_id: e.source === "material" ? e.id : null
		});
	}
	function d() {
		t.onselect?.({
			item_name: J(a).trim(),
			unit: "ea",
			unit_price: 0,
			category_type: i(),
			is_custom: !0,
			material_id: null
		});
	}
	hn(() => {
		F(o, 0);
	}), hn(() => {
		s?.focus();
	});
	var f = vi(), p = L(f);
	Gr(p), ti(p, (e) => s = e, () => s);
	var m = R(p, 2), h = (e) => {
		var t = gi();
		Er(t, 21, () => J(c), Sr, (e, t, n) => {
			var r = hi(), i = L(r), a = L(i), s = L(a, !0);
			O(a);
			var c = R(a, 2), l = L(c, !0);
			O(c), O(i);
			var d = R(i, 2), f = L(d), p = L(f);
			O(f);
			var m = R(f, 2), h = L(m);
			O(m), O(d), O(r), z((e) => {
				Lr(r, 1, `w-full text-left px-3 py-2 text-sm flex items-center justify-between
						hover:bg-slate-50 ${n === J(o) ? "bg-blue-50" : ""}`), Q(s, J(t).name), Q(l, J(t).category), Q(p, `$${e ?? ""}`), Q(h, `/ ${J(t).unit ?? ""}`);
			}, [() => J(t).price.toFixed(2)]), Y("click", r, () => u(J(t))), cr("mouseenter", r, () => F(o, n, !0)), Z(e, r);
		}), O(t), Z(e, t);
	}, g = (e) => {
		var t = _i(), n = L(t), r = R(L(n)), i = L(r);
		O(r), O(n), O(t), z(() => Q(i, `"${J(a) ?? ""}"`)), Y("click", n, d), Z(e, t);
	};
	$(m, (e) => {
		J(c).length > 0 ? e(h) : J(a).length > 0 && e(g, 1);
	}), O(f), Y("keydown", p, l), Zr(p, () => J(a), (e) => F(a, e)), Z(e, f), Re();
}
lr(["keydown", "click"]);
//#endregion
//#region node_modules/nanoid/url-alphabet/index.js
var bi = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict", xi = (e = 21) => {
	let t = "", n = crypto.getRandomValues(new Uint8Array(e |= 0));
	for (; e--;) t += bi[n[e] & 63];
	return t;
}, Si = /* @__PURE__ */ X("<div class=\"mb-2\"><!></div>"), Ci = /* @__PURE__ */ X("<table class=\"w-full\"><tbody></tbody></table>"), wi = /* @__PURE__ */ X("<div class=\"ml-4 mt-2\"><div class=\"flex items-center gap-2 mb-1\"><span class=\"text-xs font-semibold text-slate-500 uppercase tracking-wide\"> </span> <span class=\"text-xs text-slate-400\"> </span> <button class=\"text-xs text-blue-500 hover:text-blue-700 ml-auto\">+ Add Item</button></div> <!> <!></div>");
function Ti(e, t) {
	Le(t, !0);
	let n = /* @__PURE__ */ P(!1);
	function r(e) {
		t.group.line_items.push({
			id: xi(),
			category_type: e.category_type,
			item_name: e.item_name,
			quantity: 1,
			unit: e.unit,
			unit_price: e.unit_price,
			is_custom: e.is_custom,
			material_id: e.material_id,
			price_override: !1,
			description: null,
			sort_order: t.group.line_items.length,
			component_group_id: t.group.id
		}), F(n, !1), t.onchange?.();
	}
	function i(e) {
		let n = t.group.line_items.findIndex((t) => t.id === e);
		n !== -1 && (t.group.line_items.splice(n, 1), t.onchange?.());
	}
	var a = wi(), o = L(a), s = L(o), c = L(s, !0);
	O(s);
	var l = R(s, 2), u = L(l);
	O(l);
	var d = R(l, 2);
	O(o);
	var f = R(o, 2), p = (e) => {
		var i = Si();
		yi(L(i), {
			get materialsDb() {
				return t.materialsDb;
			},
			get ratesDb() {
				return t.ratesDb;
			},
			categoryType: "materials",
			onselect: r,
			oncancel: () => F(n, !1)
		}), O(i), Z(e, i);
	};
	$(f, (e) => {
		J(n) && e(p);
	});
	var m = R(f, 2), h = (e) => {
		var n = Ci(), r = L(n);
		Er(r, 21, () => t.group.line_items, (e) => e.id, (e, n) => {
			mi(e, {
				get item() {
					return J(n);
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
				ondelete: i
			});
		}), O(r), O(n), Z(e, n);
	};
	$(m, (e) => {
		t.group.line_items.length > 0 && e(h);
	}), O(a), z(() => {
		Q(c, t.group.name), Q(u, `(${t.group.line_items.length ?? ""})`);
	}), Y("click", d, () => F(n, !J(n))), Z(e, a), Re();
}
lr(["click"]);
//#endregion
//#region src/lib/SubcategoryBlock.svelte
var Ei = /* @__PURE__ */ X("<span class=\"text-xs px-1.5 py-0.5 rounded bg-amber-50 text-amber-600 font-medium\">overrides</span>"), Di = /* @__PURE__ */ X("<span class=\"text-xs px-1.5 py-0.5 rounded bg-green-50 text-green-600 font-medium\"> </span>"), Oi = /* @__PURE__ */ X("<span> </span>"), ki = /* @__PURE__ */ X("<div class=\"flex items-center gap-3 py-1.5 px-2 bg-amber-50 rounded text-xs mb-2\"><span class=\"font-medium text-amber-700\">Markup:</span> <!></div>"), Ai = /* @__PURE__ */ X("<table class=\"w-full\"><thead><tr class=\"text-xs text-slate-400 uppercase tracking-wide border-b border-slate-100\"><th class=\"px-1 py-1 text-left w-24\">Type</th><th class=\"px-1 py-1 text-left\">Name</th><th class=\"px-1 py-1 text-right w-20\">Qty</th><th class=\"px-1 py-1 text-center w-16\">Unit</th><th class=\"px-1 py-1 text-right w-24\">Price</th><th class=\"px-2 py-1 text-right w-16\">Markup</th><th class=\"px-2 py-1 text-right w-24\">w/ Markup</th><th class=\"px-2 py-1 text-right w-24\">Total</th><th class=\"w-8\"></th></tr></thead><tbody></tbody></table>"), ji = /* @__PURE__ */ X("<div class=\"mt-2\"><!></div>"), Mi = /* @__PURE__ */ X("<button class=\"mt-2 text-sm text-blue-500 hover:text-blue-700 flex items-center gap-1\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Item</button>"), Ni = /* @__PURE__ */ X("<span class=\"text-green-600\"> </span>"), Pi = /* @__PURE__ */ X("<div class=\"px-4 pb-3\"><!> <!> <!> <!> <div class=\"flex justify-between items-center mt-2 pt-2 border-t border-slate-100 text-sm\"><span class=\"text-slate-500\"> <!></span> <span class=\"font-mono font-semibold text-slate-700\"> </span></div></div>"), Fi = /* @__PURE__ */ X("<div class=\"border-t border-slate-200\"><button class=\"w-full flex items-center justify-between px-4 py-2 hover:bg-slate-50 transition-colors text-left\"><div class=\"flex items-center gap-2\"><svg fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> <span class=\"font-medium text-slate-600 text-sm\"> </span> <span class=\"text-xs text-slate-400\"> </span> <!> <!></div> <span class=\"font-mono text-sm font-semibold text-slate-700\"> </span></button> <!></div>");
function Ii(e, t) {
	Le(t, !0);
	let n = ni(t, "collapsed", 3, !1), r = ni(t, "materialsDb", 19, () => []), i = ni(t, "ratesDb", 19, () => []), a = /* @__PURE__ */ P(Wt(n())), o = /* @__PURE__ */ P(!1), s = /* @__PURE__ */ N(() => oi(t.subcat, t.globals)), c = /* @__PURE__ */ N(() => [
		"materials",
		"labor",
		"equipment",
		"subs",
		"other"
	].map((e) => ({
		type: e,
		value: ri(e, t.globals, t.subcat.markup_overrides, t.subcat.markup_enabled),
		isOverride: t.subcat.markup_overrides[e] != null,
		isDisabled: !t.subcat.markup_enabled[e]
	}))), l = /* @__PURE__ */ N(() => J(c).some((e) => e.isOverride || e.isDisabled)), u = /* @__PURE__ */ N(() => t.subcat.line_items.length + t.subcat.component_groups.reduce((e, t) => e + t.line_items.length, 0));
	function d(e) {
		t.subcat.line_items.push({
			id: xi(),
			category_type: e.category_type,
			item_name: e.item_name,
			quantity: 1,
			unit: e.unit,
			unit_price: e.unit_price,
			is_custom: e.is_custom,
			material_id: e.material_id,
			price_override: !1,
			description: null,
			sort_order: t.subcat.line_items.length,
			component_group_id: null
		}), F(o, !1), t.onchange?.();
	}
	function f(e) {
		let n = t.subcat.line_items.findIndex((t) => t.id === e);
		n !== -1 && (t.subcat.line_items.splice(n, 1), t.onchange?.());
	}
	var p = Fi(), m = L(p), h = L(m), g = L(h), _ = R(g, 2), v = L(_, !0);
	O(_);
	var y = R(_, 2), b = L(y);
	O(y);
	var x = R(y, 2), S = (e) => {
		Z(e, Ei());
	};
	$(x, (e) => {
		J(l) && e(S);
	});
	var ee = R(x, 2), te = (e) => {
		var n = Di(), r = L(n);
		O(n), z((e) => Q(r, `+${e ?? ""} lump sum`), [() => li(t.subcat.lump_sum)]), Z(e, n);
	};
	$(ee, (e) => {
		t.subcat.lump_sum > 0 && e(te);
	}), O(h);
	var C = R(h, 2), ne = L(C, !0);
	O(C), O(m);
	var re = R(m, 2), ie = (e) => {
		var n = Pi(), a = L(n), p = (e) => {
			var t = ki();
			Er(R(L(t), 2), 17, () => J(c), Sr, (e, t) => {
				var n = Oi(), r = L(n);
				O(n), z((e) => {
					Lr(n, 1, Pr(J(t).isDisabled ? "text-slate-400 line-through" : J(t).isOverride ? "text-amber-700 font-medium" : "text-slate-500")), Q(r, `${J(t).type ?? ""} ${e ?? ""}`);
				}, [() => ui(J(t).value)]), Z(e, n);
			}), O(t), Z(e, t);
		};
		$(a, (e) => {
			J(l) && e(p);
		});
		var m = R(a, 2), h = (e) => {
			var n = Ai(), r = R(L(n));
			Er(r, 21, () => t.subcat.line_items, (e) => e.id, (e, n) => {
				mi(e, {
					get item() {
						return J(n);
					},
					get globals() {
						return t.globals;
					},
					get markupOverrides() {
						return t.subcat.markup_overrides;
					},
					get markupEnabled() {
						return t.subcat.markup_enabled;
					},
					get onchange() {
						return t.onchange;
					},
					ondelete: f
				});
			}), O(r), O(n), Z(e, n);
		};
		$(m, (e) => {
			(J(u) > 0 || t.subcat.line_items.length > 0) && e(h);
		});
		var g = R(m, 2);
		Er(g, 17, () => t.subcat.component_groups, (e) => e.id, (e, n) => {
			Ti(e, {
				get group() {
					return J(n);
				},
				get globals() {
					return t.globals;
				},
				get markupOverrides() {
					return t.subcat.markup_overrides;
				},
				get markupEnabled() {
					return t.subcat.markup_enabled;
				},
				get onchange() {
					return t.onchange;
				},
				get materialsDb() {
					return r();
				},
				get ratesDb() {
					return i();
				}
			});
		});
		var _ = R(g, 2), v = (e) => {
			var t = ji();
			yi(L(t), {
				get materialsDb() {
					return r();
				},
				get ratesDb() {
					return i();
				},
				categoryType: "materials",
				onselect: d,
				oncancel: () => F(o, !1)
			}), O(t), Z(e, t);
		}, y = (e) => {
			var t = Mi();
			Y("click", t, () => F(o, !0)), Z(e, t);
		};
		$(_, (e) => {
			J(o) ? e(v) : e(y, -1);
		});
		var b = R(_, 2), x = L(b), S = L(x), ee = R(S), te = (e) => {
			var n = Ni(), r = L(n);
			O(n), z((e) => Q(r, `+ ${e ?? ""} lump sum`), [() => li(t.subcat.lump_sum)]), Z(e, n);
		};
		$(ee, (e) => {
			t.subcat.lump_sum > 0 && e(te);
		}), O(x);
		var C = R(x, 2), ne = L(C, !0);
		O(C), O(b), O(n), z((e, t) => {
			Q(S, `Subtotal: ${e ?? ""} `), Q(ne, t);
		}, [() => li(J(s).base), () => li(J(s).withMarkup)]), Z(e, n);
	};
	$(re, (e) => {
		J(a) || e(ie);
	}), O(p), z((e) => {
		Lr(g, 0, `w-4 h-4 text-slate-400 transition-transform ${J(a) ? "" : "rotate-90"}`), Q(v, t.subcat.name), Q(b, `${J(u) ?? ""} item${J(u) === 1 ? "" : "s"}`), Q(ne, e);
	}, [() => li(J(s).withMarkup)]), Y("click", m, () => F(a, !J(a))), Z(e, p), Re();
}
lr(["click"]);
//#endregion
//#region src/lib/SectionBlock.svelte
var Li = /* @__PURE__ */ X("<div class=\"px-4 py-8 text-center text-slate-400 text-sm\">No subcategories yet</div>"), Ri = /* @__PURE__ */ X("<div class=\"mb-4 border border-slate-200 rounded-lg overflow-hidden bg-white\"><button class=\"w-full flex items-center justify-between px-4 py-3 bg-slate-800 text-white hover:bg-slate-700 transition-colors text-left\"><div class=\"flex items-center gap-3\"><svg fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> <span class=\"font-semibold\"> </span> <span class=\"text-xs text-slate-400\"> </span></div> <span class=\"font-mono font-semibold\"> </span></button> <!></div>");
function zi(e, t) {
	Le(t, !0);
	let n = ni(t, "collapsed", 3, !1), r = ni(t, "materialsDb", 19, () => []), i = ni(t, "ratesDb", 19, () => []), a = /* @__PURE__ */ P(Wt(n())), o = /* @__PURE__ */ N(() => si(t.section, t.globals)), s = /* @__PURE__ */ N(() => t.section.subcategories.reduce((e, t) => e + t.line_items.length + t.component_groups.reduce((e, t) => e + t.line_items.length, 0), 0));
	var c = Ri(), l = L(c), u = L(l), d = L(u), f = R(d, 2), p = L(f, !0);
	O(f);
	var m = R(f, 2), h = L(m);
	O(m), O(u);
	var g = R(u, 2), _ = L(g, !0);
	O(g), O(l);
	var v = R(l, 2), y = (e) => {
		var n = gr(), a = en(n), o = (e) => {
			Z(e, Li());
		}, s = (e) => {
			var n = gr();
			Er(en(n), 17, () => t.section.subcategories, (e) => e.id, (e, n) => {
				Ii(e, {
					get subcat() {
						return J(n);
					},
					get globals() {
						return t.globals;
					},
					get onchange() {
						return t.onchange;
					},
					get materialsDb() {
						return r();
					},
					get ratesDb() {
						return i();
					}
				});
			}), Z(e, n);
		};
		$(a, (e) => {
			t.section.subcategories.length === 0 ? e(o) : e(s, -1);
		}), Z(e, n);
	};
	$(v, (e) => {
		J(a) || e(y);
	}), O(c), z((e) => {
		Lr(d, 0, `w-4 h-4 text-slate-400 transition-transform ${J(a) ? "" : "rotate-90"}`), Q(p, t.section.name), Q(h, `${t.section.subcategories.length ?? ""} subcategor${t.section.subcategories.length === 1 ? "y" : "ies"}
				· ${J(s) ?? ""} item${J(s) === 1 ? "" : "s"}`), Q(_, e);
	}, [() => li(J(o).withMarkup)]), Y("click", l, () => F(a, !J(a))), Z(e, c), Re();
}
lr(["click"]);
//#endregion
//#region src/lib/FooterSummary.svelte
var Bi = /* @__PURE__ */ X("<div class=\"flex flex-col\"><span class=\"text-xs text-slate-400 uppercase tracking-wide\"> </span> <span class=\"font-mono\"> </span></div>"), Vi = /* @__PURE__ */ X("<span class=\"text-slate-500 text-xs\">No items yet</span>"), Hi = /* @__PURE__ */ X("<div class=\"fixed bottom-0 left-0 right-0 bg-slate-900 text-white border-t border-slate-700 z-20\"><div class=\"flex items-center justify-between px-4 py-3\"><div class=\"flex items-center gap-6 text-sm\"><div class=\"flex flex-col\"><span class=\"text-xs text-slate-400 uppercase tracking-wide\">Base Cost</span> <span class=\"font-mono\"> </span></div> <div class=\"w-px h-8 bg-slate-700\"></div> <!> <!></div> <div class=\"flex flex-col items-end\"><span class=\"text-xs text-slate-400 uppercase tracking-wide\">Total</span> <span class=\"font-mono text-lg font-bold\"> </span></div></div></div>");
function Ui(e, t) {
	Le(t, !0);
	let n = /* @__PURE__ */ N(() => ci(t.estimate)), r = /* @__PURE__ */ N(() => Object.entries(J(n).byType).filter(([, e]) => e > 0).map(([e, t]) => ({
		type: e,
		value: t
	}))), i = {
		materials: "Materials",
		labor: "Labor",
		equipment: "Equipment",
		subs: "Subs",
		other: "Other"
	};
	var a = Hi(), o = L(a), s = L(o), c = L(s), l = R(L(c), 2), u = L(l, !0);
	O(l), O(c);
	var d = R(c, 4);
	Er(d, 17, () => J(r), Sr, (e, t) => {
		let n = () => J(t).type, r = () => J(t).value;
		var a = Bi(), o = L(a), s = L(o, !0);
		O(o);
		var c = R(o, 2), l = L(c, !0);
		O(c), O(a), z((e) => {
			Q(s, i[n()]), Q(l, e);
		}, [() => li(r())]), Z(e, a);
	});
	var f = R(d, 2), p = (e) => {
		Z(e, Vi());
	};
	$(f, (e) => {
		J(r).length === 0 && e(p);
	}), O(s);
	var m = R(s, 2), h = R(L(m), 2), g = L(h, !0);
	O(h), O(m), O(o), O(a), z((e, t) => {
		Q(u, e), Q(g, t);
	}, [() => li(J(n).base), () => li(J(n).withMarkup)]), Z(e, a), Re();
}
//#endregion
//#region src/lib/SaveStatus.svelte
var Wi = /* @__PURE__ */ X("<div><span></span> </div>");
function Gi(e, t) {
	Le(t, !0);
	let n = /* @__PURE__ */ N(() => {
		switch (t.status) {
			case "clean":
				if (t.savedAt) {
					let e = Math.round((Date.now() - t.savedAt.getTime()) / 1e3);
					return e < 5 ? "Saved just now" : e < 60 ? `Saved ${e}s ago` : `Saved ${Math.round(e / 60)}m ago`;
				}
				return "Up to date";
			case "dirty": return "Unsaved changes";
			case "saving": return "Saving...";
			case "error": return "Save failed — retrying";
			default: return "";
		}
	}), r = /* @__PURE__ */ N(() => {
		switch (t.status) {
			case "clean": return "text-green-400";
			case "dirty": return "text-amber-400";
			case "saving": return "text-blue-400";
			case "error": return "text-red-400";
			default: return "text-slate-400";
		}
	}), i = /* @__PURE__ */ N(() => {
		switch (t.status) {
			case "clean": return "bg-green-400";
			case "dirty": return "bg-amber-400";
			case "saving": return "bg-blue-400 animate-pulse";
			case "error": return "bg-red-400";
			default: return "bg-slate-400";
		}
	});
	var a = Wi(), o = L(a), s = R(o);
	O(a), z(() => {
		Lr(a, 1, `flex items-center gap-1.5 text-xs ${J(r) ?? ""}`), Lr(o, 1, `w-2 h-2 rounded-full ${J(i) ?? ""}`), Q(s, ` ${J(n) ?? ""}`);
	}), Z(e, a), Re();
}
//#endregion
//#region src/lib/autosave.svelte.js
function Ki(e, t = 2e3) {
	let n = /* @__PURE__ */ P("clean"), r = /* @__PURE__ */ P(null), i = null, a = null;
	function o(e) {
		a = e;
	}
	function s() {
		F(n, "dirty"), i && clearTimeout(i), i = setTimeout(() => c(), t);
	}
	async function c() {
		if (i &&= (clearTimeout(i), null), !a) return null;
		let o = a();
		if (!o) return null;
		F(n, "saving");
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
			return F(n, "clean"), F(r, /* @__PURE__ */ new Date(), !0), i;
		} catch (e) {
			return F(n, "error"), console.error("Auto-save failed:", e.message), i = setTimeout(() => c(), t * 2), null;
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
			return J(n);
		},
		get savedAt() {
			return J(r);
		}
	};
}
//#endregion
//#region src/EstimateBuilder.svelte
var qi = /* @__PURE__ */ X("<div class=\"flex items-center justify-center h-64\"><div class=\"text-slate-500\">Loading estimate...</div></div>"), Ji = /* @__PURE__ */ X("<div class=\"flex items-center justify-center h-64\"><div class=\"text-red-500\"> </div></div>"), Yi = /* @__PURE__ */ X("<div class=\"text-center py-16 text-slate-400\"><svg class=\"w-12 h-12 mx-auto mb-3 text-slate-300\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"1.5\" d=\"M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z\"></path></svg> <p class=\"text-lg font-medium\">No sections yet</p> <p class=\"text-sm mt-1\">Add a section to start building your estimate.</p></div>"), Xi = /* @__PURE__ */ X("<div class=\"estimate-builder pb-16\"><div class=\"sticky top-0 z-10 bg-slate-900 text-white px-4 py-3 flex items-center justify-between border-b border-slate-700\"><div><h1 class=\"text-lg font-semibold\"> </h1> <span class=\"text-xs text-slate-400\">Estimate Builder</span></div> <div class=\"flex items-center gap-3\"><!> <span class=\"text-xs px-2 py-1 rounded bg-slate-700 text-slate-300 uppercase tracking-wide\"> </span></div></div> <div class=\"bg-slate-50 border-b border-slate-200 px-4 py-2\"><div class=\"flex items-center gap-4 text-sm\"><span class=\"font-medium text-slate-600\">Global Markup:</span> <span class=\"px-2 py-0.5 rounded bg-blue-50 text-blue-700 font-mono text-xs\"> </span> <span class=\"px-2 py-0.5 rounded bg-amber-50 text-amber-700 font-mono text-xs\"> </span> <span class=\"px-2 py-0.5 rounded bg-purple-50 text-purple-700 font-mono text-xs\"> </span> <span class=\"px-2 py-0.5 rounded bg-green-50 text-green-700 font-mono text-xs\"> </span> <span class=\"px-2 py-0.5 rounded bg-slate-100 text-slate-600 font-mono text-xs\"> </span></div></div> <div class=\"p-4\"><!></div> <!></div>");
function Zi(e, t) {
	Le(t, !0);
	let n = /* @__PURE__ */ P(null), r = /* @__PURE__ */ P(null), i = /* @__PURE__ */ P(!0), a = Ki(t.projectId);
	async function o() {
		try {
			let e = await fetch(`/api/estimate/${t.projectId}`);
			if (!e.ok) throw Error(`Failed to load estimate: ${e.status}`);
			F(n, await e.json(), !0), a.register(() => J(n));
		} catch (e) {
			F(r, e.message, !0);
		} finally {
			F(i, !1);
		}
	}
	hn(() => {
		t.projectId && o();
	}), hn(() => {
		function e(e) {
			(a.status === "dirty" || a.status === "saving") && e.preventDefault();
		}
		return window.addEventListener("beforeunload", e), () => {
			window.removeEventListener("beforeunload", e), a.destroy();
		};
	});
	function s() {
		a.markDirty();
	}
	var c = gr(), l = en(c), u = (e) => {
		Z(e, qi());
	}, d = (e) => {
		var t = Ji(), n = L(t), i = L(n, !0);
		O(n), O(t), z(() => Q(i, J(r))), Z(e, t);
	}, f = (e) => {
		var t = Xi(), r = L(t), i = L(r), o = L(i), c = L(o, !0);
		O(o), Oe(2), O(i);
		var l = R(i, 2), u = L(l);
		Gi(u, {
			get status() {
				return a.status;
			},
			get savedAt() {
				return a.savedAt;
			}
		});
		var d = R(u, 2), f = L(d, !0);
		O(d), O(l), O(r);
		var p = R(r, 2), m = L(p), h = R(L(m), 2), g = L(h);
		O(h);
		var _ = R(h, 2), v = L(_);
		O(_);
		var y = R(_, 2), b = L(y);
		O(y);
		var x = R(y, 2), S = L(x);
		O(x);
		var ee = R(x, 2), te = L(ee);
		O(ee), O(m), O(p);
		var C = R(p, 2), ne = L(C), re = (e) => {
			Z(e, Yi());
		}, ie = (e) => {
			var t = gr();
			Er(en(t), 17, () => J(n).sections, (e) => e.id, (e, t) => {
				zi(e, {
					get section() {
						return J(t);
					},
					get globals() {
						return J(n).globals;
					},
					onchange: s,
					get materialsDb() {
						return J(n).materials_db;
					},
					get ratesDb() {
						return J(n).rates_db;
					}
				});
			}), Z(e, t);
		};
		$(ne, (e) => {
			J(n).sections.length === 0 ? e(re) : e(ie, -1);
		}), O(C), Ui(R(C, 2), { get estimate() {
			return J(n);
		} }), O(t), z((e, t, r, i, a) => {
			Q(c, J(n).project.name), Q(f, J(n).project.status), Q(g, `Materials ${e ?? ""}`), Q(v, `Labor ${t ?? ""}`), Q(b, `Equipment ${r ?? ""}`), Q(S, `Subs ${i ?? ""}`), Q(te, `Other ${a ?? ""}`);
		}, [
			() => ui(J(n).globals.materials_markup),
			() => ui(J(n).globals.labor_markup),
			() => ui(J(n).globals.equipment_markup),
			() => ui(J(n).globals.subs_markup),
			() => ui(J(n).globals.other_markup)
		]), Z(e, t);
	};
	$(l, (e) => {
		J(i) ? e(u) : J(r) ? e(d, 1) : J(n) && e(f, 2);
	}), Z(e, c), Re();
}
//#endregion
//#region src/main.js
var Qi = document.getElementById("estimate-root");
if (Qi) {
	let e = Qi.dataset.projectId;
	_r(Zi, {
		target: Qi,
		props: { projectId: e }
	});
}
//#endregion
