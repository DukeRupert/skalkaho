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
var m = 1024, h = 2048, g = 4096, _ = 8192, v = 16384, y = 32768, b = 1 << 25, x = 65536, S = 1 << 19, ee = 1 << 20, te = 1 << 25, C = 65536, ne = 1 << 21, re = 1 << 22, w = 1 << 23, T = Symbol("$state"), ie = Symbol("legacy props"), ae = Symbol(""), E = new class extends Error {
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
var ve = {}, D = Symbol(), ye = "http://www.w3.org/1999/xhtml";
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
var O = !1;
function Ce(e) {
	O = e;
}
var k;
function we(e) {
	if (e === null) throw be(), ve;
	return k = e;
}
function Te() {
	return we(/* @__PURE__ */ Zt(k));
}
function A(e) {
	if (O) {
		if (/* @__PURE__ */ Zt(k) !== null) throw be(), ve;
		k = e;
	}
}
function Ee(e = 1) {
	if (O) {
		for (var t = e, n = k; t--;) n = /* @__PURE__ */ Zt(n);
		k = n;
	}
}
function De(e = !0) {
	for (var t = 0, n = k;;) {
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
function Oe(e) {
	if (!e || e.nodeType !== 8) throw be(), ve;
	return e.data;
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/equality.js
function ke(e) {
	return e === this.v;
}
function Ae(e, t) {
	return e == e ? e !== t || typeof e == "object" && !!e || typeof e == "function" : t == t;
}
function je(e) {
	return !Ae(e, this.v);
}
//#endregion
//#region node_modules/svelte/src/internal/flags/index.js
var Me = !1, Ne = !1, j = null;
function Pe(e) {
	j = e;
}
function Fe(e, t = !1, n) {
	j = {
		p: j,
		i: !1,
		c: null,
		e: null,
		s: e,
		x: null,
		r: W,
		l: Ne && !t ? {
			s: null,
			u: null,
			$: []
		} : null
	};
}
function Ie(e) {
	var t = j, n = t.e;
	if (n !== null) {
		t.e = null;
		for (var r of n) hn(r);
	}
	return e !== void 0 && (t.x = e), t.i = !0, j = t.p, e ?? {};
}
function Le() {
	return !Ne || j !== null && j.l === null;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/task.js
var Re = [];
function ze() {
	var e = Re;
	Re = [], f(e);
}
function Be(e) {
	if (Re.length === 0 && !et) {
		var t = Re;
		queueMicrotask(() => {
			t === Re && ze();
		});
	}
	Re.push(e);
}
function Ve() {
	for (; Re.length > 0;) ze();
}
function He(e) {
	var t = W;
	if (t === null) return U.f |= w, e;
	if (!(t.f & 32768) && !(t.f & 4)) throw e;
	Ue(e, t);
}
function Ue(e, t) {
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
var We = ~(h | g | m);
function M(e, t) {
	e.f = e.f & We | t;
}
function Ge(e) {
	e.f & 512 || e.deps === null ? M(e, m) : M(e, g);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/utils.js
function Ke(e) {
	if (e !== null) for (let t of e) !(t.f & 2) || !(t.f & 65536) || (t.f ^= C, Ke(t.deps));
}
function qe(e, t, n) {
	e.f & 2048 ? t.add(e) : e.f & 4096 && n.add(e), Ke(e.deps), M(e, m);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/store.js
var Je = !1, Ye = !1;
function Xe(e) {
	var t = Ye;
	try {
		return Ye = !1, [e(), Ye];
	} finally {
		Ye = t;
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/batch.js
var Ze = /* @__PURE__ */ new Set(), N = null, Qe = null, P = null, $e = null, et = !1, tt = !1, nt = null, rt = null, it = 0, at = 1, ot = class e {
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
			for (var n of t.d) M(n, h), this.schedule(n);
			for (n of t.m) M(n, g), this.schedule(n);
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
		if (N = null, i.length > 0) {
			var a = e.ensure();
			for (let e of i) a.schedule(e);
		}
		if (nt = null, rt = null, this.#u()) {
			this.#p(r), this.#p(n);
			for (let [e, t] of this.#c) mt(e, t);
		} else {
			this.#n === 0 && Ze.delete(this), this.#o.clear(), this.#s.clear();
			for (let e of this.#e) e(this);
			this.#e.clear(), Qe = this, ut(r), ut(n), Qe = null, this.#i?.resolve();
		}
		var o = N;
		if (this.#a.length > 0) {
			let e = o ??= this;
			e.#a.push(...this.#a.filter((t) => !e.#a.includes(t)));
		}
		o !== null && (Ze.add(o), o.#d()), Ze.has(this) || this.#m();
	}
	#f(e, t, n) {
		e.f ^= m;
		for (var r = e.first; r !== null;) {
			var i = r.f, a = (i & 96) != 0;
			if (!(a && i & 1024 || i & 8192 || this.#c.has(r)) && r.fn !== null) {
				a ? r.f ^= m : i & 4 ? t.push(r) : Me && i & 16777224 ? n.push(r) : Yn(r) && (i & 16 && this.#s.add(r), er(r));
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
		for (var t = 0; t < e.length; t += 1) qe(e[t], this.#o, this.#s);
	}
	capture(e, t) {
		t !== D && !this.previous.has(e) && this.previous.set(e, t), e.f & 8388608 || (this.current.set(e, e.v), P?.set(e, e.v));
	}
	activate() {
		N = this;
	}
	deactivate() {
		N = null, P = null;
	}
	flush() {
		try {
			if (tt = !0, N = this, !this.#u()) {
				for (let e of this.#o) this.#s.delete(e), M(e, h), this.schedule(e);
				for (let e of this.#s) M(e, g), this.schedule(e);
			}
			this.#d();
		} finally {
			it = 0, $e = null, nt = null, rt = null, tt = !1, N = null, P = null, Pt.clear();
		}
	}
	discard() {
		for (let e of this.#t) e(this);
		this.#t.clear();
	}
	#m() {
		for (let s of Ze) {
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
		--this.#n, e && --this.#r, !(this.#l || t) && (this.#l = !0, Be(() => {
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
		if (N === null) {
			let t = N = new e();
			tt || (Ze.add(N), et || Be(() => {
				N === t && t.flush();
			}));
		}
		return N;
	}
	apply() {
		if (!Me || !this.is_fork && Ze.size === 1) {
			P = null;
			return;
		}
		P = new Map(this.current);
		for (let e of Ze) if (e !== this) for (let [t, n] of e.previous) P.has(t) || P.set(t, n);
	}
	schedule(e) {
		if ($e = e, e.b?.is_pending && e.f & 16777228 && !(e.f & 32768)) {
			e.b.defer_effect(e);
			return;
		}
		for (var t = e; t.parent !== null;) {
			t = t.parent;
			var n = t.f;
			if (nt !== null && t === W && (Me || (U === null || !(U.f & 2)) && !Je)) return;
			if (n & 96) {
				if (!(n & 1024)) return;
				t.f ^= m;
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
		for (e && (N !== null && !N.is_fork && N.flush(), n = e());;) {
			if (Ve(), N === null) return n;
			N.flush();
		}
	} finally {
		et = t;
	}
}
function ct() {
	try {
		fe();
	} catch (e) {
		Ue(e, $e);
	}
}
var lt = null;
function ut(e) {
	var t = e.length;
	if (t !== 0) {
		for (var n = 0; n < t;) {
			var r = e[n++];
			if (!(r.f & 24576) && Yn(r) && (lt = /* @__PURE__ */ new Set(), er(r), r.deps === null && r.first === null && r.nodes === null && r.teardown === null && r.ac === null && En(r), lt?.size > 0)) {
				Pt.clear();
				for (let e of lt) {
					if (e.f & 24576) continue;
					let t = [e], n = e.parent;
					for (; n !== null;) lt.has(n) && (lt.delete(n), t.push(n)), n = n.parent;
					for (let e = t.length - 1; e >= 0; e--) {
						let n = t[e];
						n.f & 24576 || er(n);
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
		e & 2 ? dt(i, t, n, r) : e & 4194320 && !(e & 2048) && ft(i, t, r) && (M(i, h), pt(i));
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
	N.schedule(e);
}
function mt(e, t) {
	if (!(e.f & 32 && e.f & 1024)) {
		e.f & 2048 ? t.d.push(e) : e.f & 4096 && t.m.push(e), M(e, m);
		for (var n = e.first; n !== null;) mt(n, t), n = n.next;
	}
}
function ht(e) {
	M(e, m);
	for (var t = e.first; t !== null;) ht(t), t = t.next;
}
//#endregion
//#region node_modules/svelte/src/reactivity/create-subscriber.js
function gt(e) {
	let t = 0, n = It(0), r;
	return () => {
		fn() && (K(n), yn(() => (t === 0 && (r = ir(() => e(() => Bt(n)))), t += 1, () => {
			Be(() => {
				--t, t === 0 && (r?.(), r = void 0, Bt(n));
			});
		})));
	};
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/boundary.js
var _t = x | S;
function vt(e, t, n, r) {
	new yt(e, t, n, r);
}
var yt = class {
	parent;
	is_pending = !1;
	transform_error;
	#e;
	#t = O ? k : null;
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
			if (O) {
				let e = this.#t;
				Te();
				let t = e.data === "[!";
				if (e.data.startsWith("[?")) {
					let t = JSON.parse(e.data.slice(2));
					this.#_(t);
				} else t ? this.#v() : this.#g();
			} else this.#y();
		}, _t), O && (this.#e = k);
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
		e && (this.is_pending = !0, this.#o = xn(() => e(this.#e)), Be(() => {
			var e = this.#c = document.createDocumentFragment(), t = Yt();
			e.append(t), this.#a = this.#x(() => xn(() => this.#r(t))), this.#u === 0 && (this.#e.before(e), this.#c = null, Dn(this.#o, () => {
				this.#o = null;
			}), this.#b(N));
		}));
	}
	#y() {
		try {
			if (this.is_pending = this.has_pending_snippet(), this.#u = 0, this.#l = 0, this.#a = xn(() => {
				this.#r(this.#e);
			}), this.#u > 0) {
				var e = this.#c = document.createDocumentFragment();
				jn(this.#a, e);
				let t = this.#n.pending;
				this.#o = xn(() => t(this.#e));
			} else this.#b(N);
		} catch (e) {
			this.error(e);
		}
	}
	#b(e) {
		this.is_pending = !1;
		for (let t of this.#f) M(t, h), e.schedule(t);
		for (let t of this.#p) M(t, g), e.schedule(t);
		this.#f.clear(), this.#p.clear();
	}
	defer_effect(e) {
		qe(e, this.#f, this.#p);
	}
	is_rendered() {
		return !this.is_pending && (!this.parent || this.parent.is_rendered());
	}
	has_pending_snippet() {
		return !!this.#n.pending;
	}
	#x(e) {
		var t = W, n = U, r = j;
		Rn(this.#i), Ln(this.#i), Pe(this.#i.ctx);
		try {
			return ot.ensure(), e();
		} catch (e) {
			return He(e), null;
		} finally {
			Rn(t), Ln(n), Pe(r);
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
		this.#S(e, t), this.#l += e, !(!this.#m || this.#d) && (this.#d = !0, Be(() => {
			this.#d = !1, this.#m && Rt(this.#m, this.#l);
		}));
	}
	get_effect_pending() {
		return this.#h(), K(this.#m);
	}
	error(e) {
		var t = this.#n.onerror;
		let n = this.#n.failed;
		if (!t && !n) throw e;
		this.#a &&= (H(this.#a), null), this.#o &&= (H(this.#o), null), this.#s &&= (H(this.#s), null), O && (we(this.#t), Ee(), we(De()));
		var r = !1, i = !1;
		let a = () => {
			if (r) {
				Se();
				return;
			}
			r = !0, i && _e(), this.#s !== null && Dn(this.#s, () => {
				this.#s = null;
			}), this.#x(() => {
				this.#y();
			});
		}, o = (e) => {
			try {
				i = !0, t?.(e, a), i = !1;
			} catch (e) {
				Ue(e, this.#i && this.#i.parent);
			}
			n && (this.#s = this.#x(() => {
				try {
					return xn(() => {
						var t = W;
						t.b = this, t.f |= 128, n(this.#e, () => e, () => a);
					});
				} catch (e) {
					return Ue(e, this.#i.parent), null;
				}
			}));
		};
		Be(() => {
			var t;
			try {
				t = this.transform_error(e);
			} catch (e) {
				Ue(e, this.#i && this.#i.parent);
				return;
			}
			typeof t == "object" && t && typeof t.then == "function" ? t.then(o, (e) => Ue(e, this.#i && this.#i.parent)) : o(t);
		});
	}
};
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/async.js
function bt(e, t, n, r) {
	let i = Le() ? wt : Et;
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
			o.f & 16384 || Ue(e, o);
		}
		St();
	}
	if (n.length === 0) {
		c.then(() => l(t.map(i)));
		return;
	}
	var u = Ct();
	function d() {
		Promise.all(n.map((e) => /* @__PURE__ */ Tt(e))).then((e) => l([...t.map(i), ...e])).catch((e) => Ue(e, o)).finally(() => u());
	}
	c ? c.then(() => {
		s(), d(), St();
	}) : d();
}
function xt() {
	var e = W, t = U, n = j, r = N;
	return function(i = !0) {
		Rn(e), Ln(t), Pe(n), i && !(e.f & 16384) && (r?.activate(), r?.apply());
	};
}
function St(e = !0) {
	Rn(null), Ln(null), Pe(null), e && N?.deactivate();
}
function Ct() {
	var e = W.b, t = N, n = e.is_rendered();
	return e.update_pending_count(1, t), t.increment(n), (r = !1) => {
		e.update_pending_count(-1, t), t.decrement(n, r);
	};
}
/* @__NO_SIDE_EFFECTS__ */
function wt(e) {
	var t = 2 | h, n = U !== null && U.f & 2 ? U : null;
	return W !== null && (W.f |= S), {
		ctx: j,
		deps: null,
		effects: null,
		equals: ke,
		f: t,
		fn: e,
		reactions: null,
		rv: 0,
		v: D,
		wv: 0,
		parent: n ?? W,
		ac: null
	};
}
/* @__NO_SIDE_EFFECTS__ */
function Tt(e, t, n) {
	let r = W;
	r === null && se();
	var i = void 0, a = It(D), o = !U, s = /* @__PURE__ */ new Map();
	return vn(() => {
		var t = W, n = p();
		i = n.promise;
		try {
			Promise.resolve(e()).then(n.resolve, n.reject).finally(St);
		} catch (e) {
			n.reject(e), St();
		}
		var c = N;
		if (o) {
			if (t.f & 32768) var l = Ct();
			if (r.b.is_rendered()) s.get(c)?.reject(E), s.delete(c);
			else {
				for (let e of s.values()) e.reject(E);
				s.clear();
			}
			s.set(c, n);
		}
		let u = (e, n = void 0) => {
			if (l && l(n === E), !(n === E || t.f & 16384)) {
				if (c.activate(), n) a.f |= w, Rt(a, n);
				else {
					a.f & 8388608 && (a.f ^= w), Rt(a, e);
					for (let [e, t] of s) {
						if (s.delete(e), e === c) break;
						t.reject(E);
					}
				}
				c.deactivate();
			}
		};
		n.promise.then(u, (e) => u(null, e || "unknown"));
	}), pn(() => {
		for (let e of s.values()) e.reject(E);
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
function F(e) {
	let t = /* @__PURE__ */ wt(e);
	return Me || Bn(t), t;
}
/* @__NO_SIDE_EFFECTS__ */
function Et(e) {
	let t = /* @__PURE__ */ wt(e);
	return t.equals = je, t;
}
function Dt(e) {
	var t = e.effects;
	if (t !== null) {
		e.effects = null;
		for (var n = 0; n < t.length; n += 1) H(t[n]);
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
	Rn(Ot(e));
	try {
		e.f &= ~C, Dt(e), t = Zn(e);
	} finally {
		Rn(n);
	}
	return t;
}
function At(e) {
	var t = kt(e);
	if (!e.equals(t) && (e.wv = Jn(), (!N?.is_fork || e.deps === null) && (e.v = t, e.deps === null))) {
		M(e, m);
		return;
	}
	Pn || (P === null ? Ge(e) : (fn() || N?.is_fork) && P.set(e, t));
}
function jt(e) {
	if (e.effects !== null) for (let t of e.effects) (t.teardown || t.ac) && (t.teardown?.(), t.ac?.abort(E), t.teardown = d, t.ac = null, $n(t, 0), Cn(t));
}
function Mt(e) {
	if (e.effects !== null) for (let t of e.effects) t.teardown && er(t);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/sources.js
var Nt = /* @__PURE__ */ new Set(), Pt = /* @__PURE__ */ new Map(), Ft = !1;
function It(e, t) {
	return {
		f: 0,
		v: e,
		reactions: null,
		equals: ke,
		rv: 0,
		wv: 0
	};
}
/* @__NO_SIDE_EFFECTS__ */
function I(e, t) {
	let n = It(e, t);
	return Bn(n), n;
}
/* @__NO_SIDE_EFFECTS__ */
function Lt(e, t = !1, n = !0) {
	let r = It(e);
	return t || (r.equals = je), Ne && n && j !== null && j.l !== null && (j.l.s ??= []).push(r), r;
}
function L(e, t, r = !1) {
	return U !== null && (!In || U.f & 131072) && Le() && U.f & 4325394 && (zn === null || !n.call(zn, e)) && ge(), Rt(e, r ? R(t) : t, rt);
}
function Rt(e, t, n = null) {
	if (!e.equals(t)) {
		var r = e.v;
		Pn ? Pt.set(e, t) : Pt.set(e, r), e.v = t;
		var i = ot.ensure();
		if (i.capture(e, r), e.f & 2) {
			let t = e;
			e.f & 2048 && kt(t), Ge(t);
		}
		e.wv = Jn(), Vt(e, h, n), Le() && W !== null && W.f & 1024 && !(W.f & 96) && (Hn === null ? Un([e]) : Hn.push(e)), !i.is_fork && Nt.size > 0 && !Ft && zt();
	}
	return t;
}
function zt() {
	Ft = !1;
	for (let e of Nt) e.f & 1024 && M(e, g), Yn(e) && er(e);
	Nt.clear();
}
function Bt(e) {
	L(e, e.v + 1);
}
function Vt(e, t, n) {
	var r = e.reactions;
	if (r !== null) for (var i = Le(), a = r.length, o = 0; o < a; o++) {
		var s = r[o], c = s.f;
		if (!(!i && s === W)) {
			var l = (c & h) === 0;
			if (l && M(s, t), c & 2) {
				var u = s;
				P?.delete(u), c & 65536 || (c & 512 && (s.f |= C), Vt(u, g, n));
			} else if (l) {
				var d = s;
				c & 16 && lt !== null && lt.add(d), n === null ? pt(d) : n.push(d);
			}
		}
	}
}
function R(t) {
	if (typeof t != "object" || !t || T in t) return t;
	let n = l(t);
	if (n !== s && n !== c) return t;
	var r = /* @__PURE__ */ new Map(), i = e(t), o = /* @__PURE__ */ I(0), u = null, d = Kn, f = (e) => {
		if (Kn === d) return e();
		var t = U, n = Kn;
		Ln(null), qn(d);
		var r = e();
		return Ln(t), qn(n), r;
	};
	return i && r.set("length", /* @__PURE__ */ I(t.length, u)), new Proxy(t, {
		defineProperty(e, t, n) {
			(!("value" in n) || n.configurable === !1 || n.enumerable === !1 || n.writable === !1) && me();
			var i = r.get(t);
			return i === void 0 ? f(() => {
				var e = /* @__PURE__ */ I(n.value, u);
				return r.set(t, e), e;
			}) : L(i, n.value, !0), !0;
		},
		deleteProperty(e, t) {
			var n = r.get(t);
			if (n === void 0) {
				if (t in e) {
					let e = f(() => /* @__PURE__ */ I(D, u));
					r.set(t, e), Bt(o);
				}
			} else L(n, D), Bt(o);
			return !0;
		},
		get(e, n, i) {
			if (n === T) return t;
			var o = r.get(n), s = n in e;
			if (o === void 0 && (!s || a(e, n)?.writable) && (o = f(() => /* @__PURE__ */ I(R(s ? e[n] : D), u)), r.set(n, o)), o !== void 0) {
				var c = K(o);
				return c === D ? void 0 : c;
			}
			return Reflect.get(e, n, i);
		},
		getOwnPropertyDescriptor(e, t) {
			var n = Reflect.getOwnPropertyDescriptor(e, t);
			if (n && "value" in n) {
				var i = r.get(t);
				i && (n.value = K(i));
			} else if (n === void 0) {
				var a = r.get(t), o = a?.v;
				if (a !== void 0 && o !== D) return {
					enumerable: !0,
					configurable: !0,
					value: o,
					writable: !0
				};
			}
			return n;
		},
		has(e, t) {
			if (t === T) return !0;
			var n = r.get(t), i = n !== void 0 && n.v !== D || Reflect.has(e, t);
			return (n !== void 0 || W !== null && (!i || a(e, t)?.writable)) && (n === void 0 && (n = f(() => /* @__PURE__ */ I(i ? R(e[t]) : D, u)), r.set(t, n)), K(n) === D) ? !1 : i;
		},
		set(e, t, n, s) {
			var c = r.get(t), l = t in e;
			if (i && t === "length") for (var d = n; d < c.v; d += 1) {
				var p = r.get(d + "");
				p === void 0 ? d in e && (p = f(() => /* @__PURE__ */ I(D, u)), r.set(d + "", p)) : L(p, D);
			}
			if (c === void 0) (!l || a(e, t)?.writable) && (c = f(() => /* @__PURE__ */ I(void 0, u)), L(c, R(n)), r.set(t, c));
			else {
				l = c.v !== D;
				var m = f(() => R(n));
				L(c, m);
			}
			var h = Reflect.getOwnPropertyDescriptor(e, t);
			if (h?.set && h.set.call(s, n), !l) {
				if (i && typeof t == "string") {
					var g = r.get("length"), _ = Number(t);
					Number.isInteger(_) && _ >= g.v && L(g, _ + 1);
				}
				Bt(o);
			}
			return !0;
		},
		ownKeys(e) {
			K(o);
			var t = Reflect.ownKeys(e).filter((e) => {
				var t = r.get(e);
				return t === void 0 || t.v !== D;
			});
			for (var [n, i] of r) i.v !== D && !(n in e) && t.push(n);
			return t;
		},
		setPrototypeOf() {
			he();
		}
	});
}
function Ht(e) {
	try {
		if (typeof e == "object" && e && T in e) return e[T];
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
function z(e, t) {
	if (!O) return /* @__PURE__ */ Xt(e);
	var n = /* @__PURE__ */ Xt(k);
	if (n === null) n = k.appendChild(Yt());
	else if (t && n.nodeType !== 3) {
		var r = Yt();
		return n?.before(r), we(r), r;
	}
	return t && nn(n), we(n), n;
}
function Qt(e, t = !1) {
	if (!O) {
		var n = /* @__PURE__ */ Xt(e);
		return n instanceof Comment && n.data === "" ? /* @__PURE__ */ Zt(n) : n;
	}
	if (t) {
		if (k?.nodeType !== 3) {
			var r = Yt();
			return k?.before(r), we(r), r;
		}
		nn(k);
	}
	return k;
}
function B(e, t = 1, n = !1) {
	let r = O ? k : e;
	for (var i; t--;) i = r, r = /* @__PURE__ */ Zt(r);
	if (!O) return r;
	if (n) {
		if (r?.nodeType !== 3) {
			var a = Yt();
			return r === null ? i?.after(a) : r.before(a), we(a), a;
		}
		nn(r);
	}
	return we(r), r;
}
function $t(e) {
	e.textContent = "";
}
function en() {
	return !Me || lt !== null ? !1 : (W.f & y) !== 0;
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
		e.autofocus = !0, Be(() => {
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
	Ln(null), Rn(null);
	try {
		return e();
	} finally {
		Ln(t), Rn(n);
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
	W === null && (U === null && de(e), ue()), Pn && le(e);
}
function un(e, t) {
	var n = t.last;
	n === null ? t.last = t.first = e : (n.next = e, e.prev = n, t.last = e);
}
function dn(e, t) {
	var n = W;
	n !== null && n.f & 8192 && (e |= _);
	var r = {
		ctx: j,
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
	if (e & 4) nt === null ? ot.ensure().schedule(r) : nt.push(r);
	else if (t !== null) {
		try {
			er(r);
		} catch (e) {
			throw H(r), e;
		}
		i.deps === null && i.teardown === null && i.nodes === null && i.first === i.last && !(i.f & 524288) && (i = i.first, e & 16 && e & 65536 && i !== null && (i.f |= x));
	}
	if (i !== null && (i.parent = n, n !== null && un(i, n), U !== null && U.f & 2 && !(e & 64))) {
		var a = U;
		(a.effects ??= []).push(i);
	}
	return r;
}
function fn() {
	return U !== null && !In;
}
function pn(e) {
	let t = dn(8, null);
	return M(t, m), t.teardown = e, t;
}
function mn(e) {
	ln("$effect");
	var t = W.f;
	if (!U && t & 32 && !(t & 32768)) {
		var n = j;
		(n.e ??= []).push(e);
	} else return hn(e);
}
function hn(e) {
	return dn(4 | ee, e);
}
function gn(e) {
	ot.ensure();
	let t = dn(64 | S, e);
	return (e = {}) => new Promise((n) => {
		e.outro ? Dn(t, () => {
			H(t), n(void 0);
		}) : (H(t), n(void 0));
	});
}
function _n(e) {
	return dn(4, e);
}
function vn(e) {
	return dn(re | S, e);
}
function yn(e, t = 0) {
	return dn(8 | t, e);
}
function V(e, t = [], n = [], r = []) {
	bt(r, t, n, (t) => {
		dn(8, () => e(...t.map(K)));
	});
}
function bn(e, t = 0) {
	return dn(16 | t, e);
}
function xn(e) {
	return dn(32 | S, e);
}
function Sn(e) {
	var t = e.teardown;
	if (t !== null) {
		let e = Pn, n = U;
		Fn(!0), Ln(null);
		try {
			t.call(null);
		} finally {
			Fn(e), Ln(n);
		}
	}
}
function Cn(e, t = !1) {
	var n = e.first;
	for (e.first = e.last = null; n !== null;) {
		let e = n.ac;
		e !== null && sn(() => {
			e.abort(E);
		});
		var r = n.next;
		n.f & 64 ? n.parent = null : H(n, t), n = r;
	}
}
function wn(e) {
	for (var t = e.first; t !== null;) {
		var n = t.next;
		t.f & 32 || H(t), t = n;
	}
}
function H(e, t = !0) {
	var n = !1;
	(t || e.f & 262144) && e.nodes !== null && e.nodes.end !== null && (Tn(e.nodes.start, e.nodes.end), n = !0), M(e, b), Cn(e, t && !n), $n(e, 0);
	var r = e.nodes && e.nodes.t;
	if (r !== null) for (let e of r) e.stop();
	Sn(e), e.f ^= b, e.f |= v;
	var i = e.parent;
	i !== null && i.first !== null && En(e), e.next = e.prev = e.teardown = e.ctx = e.deps = e.fn = e.nodes = e.ac = null;
}
function Tn(e, t) {
	for (; e !== null;) {
		var n = e === t ? null : /* @__PURE__ */ Zt(e);
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
		n && H(e), t && t();
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
		e.f ^= _, e.f & 1024 || (M(e, h), ot.ensure().schedule(e));
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
		var i = n === r ? null : /* @__PURE__ */ Zt(n);
		t.append(n), n = i;
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/legacy.js
var Mn = null, Nn = !1, Pn = !1;
function Fn(e) {
	Pn = e;
}
var U = null, In = !1;
function Ln(e) {
	U = e;
}
var W = null;
function Rn(e) {
	W = e;
}
var zn = null;
function Bn(e) {
	U !== null && (!Me || U.f & 2) && (zn === null ? zn = [e] : zn.push(e));
}
var G = null, Vn = 0, Hn = null;
function Un(e) {
	Hn = e;
}
var Wn = 1, Gn = 0, Kn = Gn;
function qn(e) {
	Kn = e;
}
function Jn() {
	return ++Wn;
}
function Yn(e) {
	var t = e.f;
	if (t & 2048) return !0;
	if (t & 2 && (e.f &= ~C), t & 4096) {
		for (var n = e.deps, r = n.length, i = 0; i < r; i++) {
			var a = n[i];
			if (Yn(a) && At(a), a.wv > e.wv) return !0;
		}
		t & 512 && P === null && M(e, m);
	}
	return !1;
}
function Xn(e, t, r = !0) {
	var i = e.reactions;
	if (i !== null && !(!Me && zn !== null && n.call(zn, e))) for (var a = 0; a < i.length; a++) {
		var o = i[a];
		o.f & 2 ? Xn(o, t, !1) : t === o && (r ? M(o, h) : o.f & 1024 && M(o, g), pt(o));
	}
}
function Zn(e) {
	var t = G, n = Vn, r = Hn, i = U, a = zn, o = j, s = In, c = Kn, l = e.f;
	G = null, Vn = 0, Hn = null, U = l & 96 ? null : e, zn = null, Pe(e.ctx), In = !1, Kn = ++Gn, e.ac !== null && (sn(() => {
		e.ac.abort(E);
	}), e.ac = null);
	try {
		e.f |= ne;
		var u = e.fn, d = u();
		e.f |= y;
		var f = e.deps, p = N?.is_fork;
		if (G !== null) {
			var m;
			if (p || $n(e, Vn), f !== null && Vn > 0) for (f.length = Vn + G.length, m = 0; m < G.length; m++) f[Vn + m] = G[m];
			else e.deps = f = G;
			if (fn() && e.f & 512) for (m = Vn; m < f.length; m++) (f[m].reactions ??= []).push(e);
		} else !p && f !== null && Vn < f.length && ($n(e, Vn), f.length = Vn);
		if (Le() && Hn !== null && !In && f !== null && !(e.f & 6146)) for (m = 0; m < Hn.length; m++) Xn(Hn[m], e);
		if (i !== null && i !== e) {
			if (Gn++, i.deps !== null) for (let e = 0; e < n; e += 1) i.deps[e].rv = Gn;
			if (t !== null) for (let e of t) e.rv = Gn;
			Hn !== null && (r === null ? r = Hn : r.push(...Hn));
		}
		return e.f & 8388608 && (e.f ^= w), d;
	} catch (e) {
		return He(e);
	} finally {
		e.f ^= ne, G = t, Vn = n, Hn = r, U = i, zn = a, Pe(o), In = s, Kn = c;
	}
}
function Qn(e, r) {
	let i = r.reactions;
	if (i !== null) {
		var a = t.call(i, e);
		if (a !== -1) {
			var o = i.length - 1;
			o === 0 ? i = r.reactions = null : (i[a] = i[o], i.pop());
		}
	}
	if (i === null && r.f & 2 && (G === null || !n.call(G, r))) {
		var s = r;
		s.f & 512 && (s.f ^= 512, s.f &= ~C), Ge(s), jt(s), $n(s, 0);
	}
}
function $n(e, t) {
	var n = e.deps;
	if (n !== null) for (var r = t; r < n.length; r++) Qn(e, n[r]);
}
function er(e) {
	var t = e.f;
	if (!(t & 16384)) {
		M(e, m);
		var n = W, r = Nn;
		W = e, Nn = !0;
		try {
			t & 16777232 ? wn(e) : Cn(e), Sn(e);
			var i = Zn(e);
			e.teardown = typeof i == "function" ? i : null, e.wv = Wn;
		} finally {
			Nn = r, W = n;
		}
	}
}
async function tr() {
	if (Me) return new Promise((e) => {
		requestAnimationFrame(() => e()), setTimeout(() => e());
	});
	await Promise.resolve(), st();
}
function K(e) {
	var t = (e.f & 2) != 0;
	if (Mn?.add(e), U !== null && !In && !(W !== null && W.f & 16384) && (zn === null || !n.call(zn, e))) {
		var r = U.deps;
		if (U.f & 2097152) e.rv < Gn && (e.rv = Gn, G === null && r !== null && r[Vn] === e ? Vn++ : G === null ? G = [e] : G.push(e));
		else {
			(U.deps ??= []).push(e);
			var i = e.reactions;
			i === null ? e.reactions = [U] : n.call(i, U) || i.push(U);
		}
	}
	if (Pn && Pt.has(e)) return Pt.get(e);
	if (t) {
		var a = e;
		if (Pn) {
			var o = a.v;
			return (!(a.f & 1024) && a.reactions !== null || rr(a)) && (o = kt(a)), Pt.set(a, o), o;
		}
		var s = (a.f & 512) == 0 && !In && U !== null && (Nn || (U.f & 512) != 0), c = (a.f & y) === 0;
		Yn(a) && (s && (a.f |= 512), At(a)), s && !c && (Mt(a), nr(a));
	}
	if (P?.has(e)) return P.get(e);
	if (e.f & 8388608) throw e.v;
	return e.v;
}
function nr(e) {
	if (e.f |= 512, e.deps !== null) for (let t of e.deps) (t.reactions ??= []).push(e), t.f & 2 && !(t.f & 512) && (Mt(t), nr(t));
}
function rr(e) {
	if (e.v === D) return !0;
	if (e.deps === null) return !1;
	for (let t of e.deps) if (Pt.has(t) || t.f & 2 && rr(t)) return !0;
	return !1;
}
function ir(e) {
	var t = In;
	try {
		return In = !0, e();
	} finally {
		In = t;
	}
}
[.../* @__PURE__ */ "allowfullscreen.async.autofocus.autoplay.checked.controls.default.disabled.formnovalidate.indeterminate.inert.ismap.loop.multiple.muted.nomodule.novalidate.open.playsinline.readonly.required.reversed.seamless.selected.webkitdirectory.defer.disablepictureinpicture.disableremoteplayback".split(".")];
var ar = ["touchstart", "touchmove"];
function or(e) {
	return ar.includes(e);
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/events.js
var sr = Symbol("events"), cr = /* @__PURE__ */ new Set(), lr = /* @__PURE__ */ new Set();
function ur(e, t, n, r = {}) {
	function i(e) {
		if (r.capture || mr.call(t, e), !e.cancelBubble) return sn(() => n?.call(this, e));
	}
	return e.startsWith("pointer") || e.startsWith("touch") || e === "wheel" ? Be(() => {
		t.addEventListener(e, i, r);
	}) : t.addEventListener(e, i, r), i;
}
function dr(e, t, n, r, i) {
	var a = {
		capture: r,
		passive: i
	}, o = ur(e, t, n, a);
	(t === document.body || t === window || t === document || t instanceof HTMLMediaElement) && pn(() => {
		t.removeEventListener(e, o, a);
	});
}
function q(e, t, n) {
	(t[sr] ??= {})[e] = n;
}
function fr(e) {
	for (var t = 0; t < e.length; t++) cr.add(e[t]);
	for (var n of lr) n(e);
}
var pr = null;
function mr(e) {
	var t = this, n = t.ownerDocument, r = e.type, a = e.composedPath?.() || [], o = a[0] || e.target;
	pr = e;
	var s = 0, c = pr === e && e[sr];
	if (c) {
		var l = a.indexOf(c);
		if (l !== -1 && (t === document || t === window)) {
			e[sr] = t;
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
		Ln(null), Rn(null);
		try {
			for (var p, m = []; o !== null;) {
				var h = o.assignedSlot || o.parentNode || o.host || null;
				try {
					var g = o[sr]?.[r];
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
			e[sr] = t, delete e.currentTarget, Ln(d), Rn(f);
		}
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/reconciler.js
var hr = globalThis?.window?.trustedTypes && /* @__PURE__ */ globalThis.window.trustedTypes.createPolicy("svelte-trusted-html", { createHTML: (e) => e });
function gr(e) {
	return hr?.createHTML(e) ?? e;
}
function _r(e) {
	var t = tn("template");
	return t.innerHTML = gr(e.replaceAll("<!>", "<!---->")), t.content;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/template.js
function vr(e, t) {
	var n = W;
	n.nodes === null && (n.nodes = {
		start: e,
		end: t,
		a: null,
		t: null
	});
}
/* @__NO_SIDE_EFFECTS__ */
function J(e, t) {
	var n = (t & 1) != 0, r = (t & 2) != 0, i, a = !e.startsWith("<!>");
	return () => {
		if (O) return vr(k, null), k;
		i === void 0 && (i = _r(a ? e : "<!>" + e), n || (i = /* @__PURE__ */ Xt(i)));
		var t = r || Gt ? document.importNode(i, !0) : i.cloneNode(!0);
		if (n) {
			var o = /* @__PURE__ */ Xt(t), s = t.lastChild;
			vr(o, s);
		} else vr(t, t);
		return t;
	};
}
function yr() {
	if (O) return vr(k, null), k;
	var e = document.createDocumentFragment(), t = document.createComment(""), n = Yt();
	return e.append(t, n), vr(t, n), e;
}
function Y(e, t) {
	if (O) {
		var n = W;
		(!(n.f & 32768) || n.nodes.end === null) && (n.nodes.end = k), Te();
		return;
	}
	e !== null && e.before(t);
}
function X(e, t) {
	var n = t == null ? "" : typeof t == "object" ? `${t}` : t;
	n !== (e.__t ??= e.nodeValue) && (e.__t = n, e.nodeValue = `${n}`);
}
function br(e, t) {
	return Sr(e, t);
}
var xr = /* @__PURE__ */ new Map();
function Sr(e, { target: t, anchor: n, props: i = {}, events: a, context: o, intro: s = !0, transformError: c }) {
	Jt();
	var l = void 0, u = gn(() => {
		var s = n ?? t.appendChild(Yt());
		vt(s, { pending: () => {} }, (t) => {
			Fe({});
			var n = j;
			if (o && (n.c = o), a && (i.$$events = a), O && vr(t, null), l = e(t, i) || {}, O && (W.nodes.end = k, k === null || k.nodeType !== 8 || k.data !== "]")) throw be(), ve;
			Ie();
		}, c);
		var u = /* @__PURE__ */ new Set(), d = (e) => {
			for (var n = 0; n < e.length; n++) {
				var r = e[n];
				if (!u.has(r)) {
					u.add(r);
					var i = or(r);
					for (let e of [t, document]) {
						var a = xr.get(e);
						a === void 0 && (a = /* @__PURE__ */ new Map(), xr.set(e, a));
						var o = a.get(r);
						o === void 0 ? (e.addEventListener(r, mr, { passive: i }), a.set(r, 1)) : a.set(r, o + 1);
					}
				}
			}
		};
		return d(r(cr)), lr.add(d), () => {
			for (var e of u) for (let n of [t, document]) {
				var r = xr.get(n), i = r.get(e);
				--i == 0 ? (n.removeEventListener(e, mr), r.delete(e), r.size === 0 && xr.delete(n)) : r.set(e, i);
			}
			lr.delete(d), s !== n && s.parentNode?.removeChild(s);
		};
	});
	return Cr.set(l, u), l;
}
var Cr = /* @__PURE__ */ new WeakMap(), wr = class {
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
				r && (H(r.effect), this.#n.delete(n));
			}
			for (let [e, r] of this.#t) {
				if (e === t || this.#r.has(e)) continue;
				let i = () => {
					if (Array.from(this.#e.values()).includes(e)) {
						var t = document.createDocumentFragment();
						jn(r, t), t.append(Yt()), this.#n.set(e, {
							effect: r,
							fragment: t
						});
					} else H(r);
					this.#r.delete(e), this.#t.delete(e);
				};
				this.#i || !n ? (this.#r.add(e), Dn(r, i, !1)) : i();
			}
		}
	};
	#o = (e) => {
		this.#e.delete(e);
		let t = Array.from(this.#e.values());
		for (let [e, n] of this.#n) t.includes(e) || (H(n.effect), this.#n.delete(e));
	};
	ensure(e, t) {
		var n = N, r = en();
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
		} else O && (this.anchor = k), this.#a(n);
	}
};
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/if.js
function Z(e, t, n = !1) {
	var r;
	O && (r = k, Te());
	var i = new wr(e), a = n ? x : 0;
	function o(e, t) {
		if (O) {
			var n = Oe(r);
			if (e !== parseInt(n.substring(1))) {
				var a = De();
				we(a), i.anchor = a, Ce(!1), i.ensure(e, t), Ce(!0);
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
function Tr(e, t) {
	return t;
}
function Er(e, t, n) {
	for (var i = [], a = t.length, o, s = t.length, c = 0; c < a; c++) {
		let n = t[c];
		Dn(n, () => {
			if (o) {
				if (o.pending.delete(n), o.done.add(n), o.pending.size === 0) {
					var t = e.outrogroups;
					Dr(e, r(o.done)), t.delete(o), t.size === 0 && (e.outrogroups = null);
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
		Dr(e, t, !l);
	} else o = {
		pending: new Set(t),
		done: /* @__PURE__ */ new Set()
	}, (e.outrogroups ??= /* @__PURE__ */ new Set()).add(o);
}
function Dr(e, t, n = !0) {
	var r;
	if (e.pending.size > 0) {
		r = /* @__PURE__ */ new Set();
		for (let t of e.pending.values()) for (let n of t) r.add(e.items.get(n).e);
	}
	for (var i = 0; i < t.length; i++) {
		var a = t[i];
		r?.has(a) ? (a.f |= te, jn(a, document.createDocumentFragment())) : H(t[i], n);
	}
}
var Or;
function kr(t, n, i, a, o, s = null) {
	var c = t, l = /* @__PURE__ */ new Map();
	if (n & 4) {
		var u = t;
		c = O ? we(/* @__PURE__ */ Xt(u)) : u.appendChild(Yt());
	}
	O && Te();
	var d = null, f = /* @__PURE__ */ Et(() => {
		var t = i();
		return e(t) ? t : t == null ? [] : r(t);
	}), p, m = /* @__PURE__ */ new Map(), h = !0;
	function g(e) {
		v.effect.f & 16384 || (v.pending.delete(e), v.fallback = d, jr(v, p, c, n, a), d !== null && (p.length === 0 ? d.f & 33554432 ? (d.f ^= te, Nr(d, null, c)) : kn(d) : Dn(d, () => {
			d = null;
		})));
	}
	function _(e) {
		v.pending.delete(e);
	}
	var v = {
		effect: bn(() => {
			p = K(f);
			var e = p.length;
			let t = !1;
			O && Oe(c) === "[!" != (e === 0) && (c = De(), we(c), Ce(!1), t = !0);
			for (var r = /* @__PURE__ */ new Set(), u = N, v = en(), y = 0; y < e; y += 1) {
				O && k.nodeType === 8 && k.data === "]" && (c = k, t = !0, Ce(!1));
				var b = p[y], x = a(b, y), S = h ? null : l.get(x);
				S ? (S.v && Rt(S.v, b), S.i && Rt(S.i, y), v && u.unskip_effect(S.e)) : (S = Mr(l, h ? c : Or ??= Yt(), b, x, y, o, n, i), h || (S.e.f |= te), l.set(x, S)), r.add(x);
			}
			if (e === 0 && s && !d && (h ? d = xn(() => s(c)) : (d = xn(() => s(Or ??= Yt())), d.f |= te)), e > r.size && ce("", "", ""), O && e > 0 && we(De()), !h) if (m.set(u, r), v) {
				for (let [e, t] of l) r.has(e) || u.skip_effect(t.e);
				u.oncommit(g), u.ondiscard(_);
			} else g(u);
			t && Ce(!0), K(f);
		}),
		flags: n,
		items: l,
		pending: m,
		outrogroups: null,
		fallback: d
	};
	h = !1, O && (c = k);
}
function Ar(e) {
	for (; e !== null && !(e.f & 32);) e = e.next;
	return e;
}
function jr(e, t, n, i, a) {
	var o = (i & 8) != 0, s = t.length, c = e.items, l = Ar(e.effect.first), u, d = null, f, p = [], m = [], h, g, _, v;
	if (o) for (v = 0; v < s; v += 1) h = t[v], g = a(h, v), _ = c.get(g).e, _.f & 33554432 || (_.nodes?.a?.measure(), (f ??= /* @__PURE__ */ new Set()).add(_));
	for (v = 0; v < s; v += 1) {
		if (h = t[v], g = a(h, v), _ = c.get(g).e, e.outrogroups !== null) for (let t of e.outrogroups) t.pending.delete(_), t.done.delete(_);
		if (_.f & 33554432) if (_.f ^= te, _ === l) Nr(_, null, n);
		else {
			var y = d ? d.next : l;
			_ === e.effect.last && (e.effect.last = _.prev), _.prev && (_.prev.next = _.next), _.next && (_.next.prev = _.prev), Pr(e, d, _), Pr(e, _, y), Nr(_, y, n), d = _, p = [], m = [], l = Ar(d.next);
			continue;
		}
		if (_.f & 8192 && (kn(_), o && (_.nodes?.a?.unfix(), (f ??= /* @__PURE__ */ new Set()).delete(_))), _ !== l) {
			if (u !== void 0 && u.has(_)) {
				if (p.length < m.length) {
					var b = m[0], x;
					d = b.prev;
					var S = p[0], ee = p[p.length - 1];
					for (x = 0; x < p.length; x += 1) Nr(p[x], b, n);
					for (x = 0; x < m.length; x += 1) u.delete(m[x]);
					Pr(e, S.prev, ee.next), Pr(e, d, S), Pr(e, ee, b), l = b, d = ee, --v, p = [], m = [];
				} else u.delete(_), Nr(_, l, n), Pr(e, _.prev, _.next), Pr(e, _, d === null ? e.effect.first : d.next), Pr(e, d, _), d = _;
				continue;
			}
			for (p = [], m = []; l !== null && l !== _;) (u ??= /* @__PURE__ */ new Set()).add(l), m.push(l), l = Ar(l.next);
			if (l === null) continue;
		}
		_.f & 33554432 || p.push(_), d = _, l = Ar(_.next);
	}
	if (e.outrogroups !== null) {
		for (let t of e.outrogroups) t.pending.size === 0 && (Dr(e, r(t.done)), e.outrogroups?.delete(t));
		e.outrogroups.size === 0 && (e.outrogroups = null);
	}
	if (l !== null || u !== void 0) {
		var C = [];
		if (u !== void 0) for (_ of u) _.f & 8192 || C.push(_);
		for (; l !== null;) !(l.f & 8192) && l !== e.fallback && C.push(l), l = Ar(l.next);
		var ne = C.length;
		if (ne > 0) {
			var re = i & 4 && s === 0 ? n : null;
			if (o) {
				for (v = 0; v < ne; v += 1) C[v].nodes?.a?.measure();
				for (v = 0; v < ne; v += 1) C[v].nodes?.a?.fix();
			}
			Er(e, C, re);
		}
	}
	o && Be(() => {
		if (f !== void 0) for (_ of f) _.nodes?.a?.apply();
	});
}
function Mr(e, t, n, r, i, a, o, s) {
	var c = o & 1 ? o & 16 ? It(n) : /* @__PURE__ */ Lt(n, !1, !1) : null, l = o & 2 ? It(i) : null;
	return {
		v: c,
		i: l,
		e: xn(() => (a(t, c ?? n, l ?? i, s), () => {
			e.delete(r);
		}))
	};
}
function Nr(e, t, n) {
	if (e.nodes) for (var r = e.nodes.start, i = e.nodes.end, a = t && !(t.f & 33554432) ? t.nodes.start : n; r !== null;) {
		var o = /* @__PURE__ */ Zt(r);
		if (a.before(r), r === i) return;
		r = o;
	}
}
function Pr(e, t, n) {
	t === null ? e.effect.first = n : t.next = n, n === null ? e.effect.last = t : n.prev = t;
}
//#endregion
//#region node_modules/svelte/src/internal/shared/attributes.js
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
	if (O || o !== n || o === void 0) {
		var s = Ir(n, r, a);
		(!O || s !== e.getAttribute("class")) && (s == null ? e.removeAttribute("class") : t ? e.className = s : e.setAttribute("class", s)), e.__className = n;
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
		if (!e(n)) return xe();
		for (var i of t.options) i.selected = n.includes(Br(i));
		return;
	}
	for (i of t.options) if (Ut(Br(i), n)) {
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
	}), pn(() => {
		t.disconnect();
	});
}
function Br(e) {
	return "__value" in e ? e.__value : e.value;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/attributes.js
var Vr = Symbol("is custom element"), Hr = Symbol("is html"), Ur = oe ? "link" : "LINK", Wr = oe ? "progress" : "PROGRESS";
function Q(e) {
	if (O) {
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
		e.__on_r = n, Be(n), on();
	}
}
function Gr(e, t) {
	var n = Jr(e);
	n.value === (n.value = t ?? void 0) || e.value === t && (t !== 0 || e.nodeName !== Wr) || (e.value = t ?? "");
}
function Kr(e, t) {
	var n = Jr(e);
	n.checked !== (n.checked = t ?? void 0) && (e.checked = t);
}
function qr(e, t, n, r) {
	var i = Jr(e);
	O && (i[t] = e.getAttribute(t), t === "src" || t === "srcset" || t === "href" && e.nodeName === Ur) || i[t] !== (i[t] = n) && (t === "loading" && (e[ae] = n), n == null ? e.removeAttribute(t) : typeof n != "string" && Xr(e).includes(t) ? e[t] = n : e.setAttribute(t, n));
}
function Jr(e) {
	return e.__attributes ??= {
		[Vr]: e.nodeName.includes("-"),
		[Hr]: e.namespaceURI === ye
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
	cn(e, "input", async (i) => {
		var a = i ? e.defaultValue : e.value;
		if (a = Qr(e) ? $r(a) : a, n(a), N !== null && r.add(N), await tr(), a !== (a = t())) {
			var o = e.selectionStart, s = e.selectionEnd, c = e.value.length;
			if (e.value = a ?? "", s !== null) {
				var l = e.value.length;
				o === s && s === c && l > c ? (e.selectionStart = l, e.selectionEnd = l) : (e.selectionStart = o, e.selectionEnd = Math.min(s, l));
			}
		}
	}), (O && e.defaultValue !== e.value || ir(t) == null && e.value) && (n(Qr(e) ? $r(e.value) : e.value), N !== null && r.add(N)), yn(() => {
		var n = t();
		if (e === document.activeElement) {
			var i = Me ? Qe : N;
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
	return e === t || e?.[T] === t;
}
function ti(e = {}, t, n, r) {
	var i = j.r, a = W;
	return _n(() => {
		var o, s;
		return yn(() => {
			o = s, s = r?.() || [], ir(() => {
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
function $(e, t, n, r) {
	var i = !Ne || (n & 2) != 0, o = (n & 8) != 0, s = (n & 16) != 0, c = r, l = !0, u = () => (l && (l = !1, c = s ? ir(r) : r), c);
	let d;
	if (o) {
		var f = T in e || ie in e;
		d = a(e, t)?.set ?? (f && t in e ? (n) => e[t] = n : void 0);
	}
	var p, m = !1;
	o ? [p, m] = Xe(() => e[t]) : p = e[t], p === void 0 && r !== void 0 && (p = u(), d && (i && pe(t), d(p)));
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
	o && K(v);
	var y = W;
	return (function(e, t) {
		if (arguments.length > 0) {
			let n = t ? K(v) : i && o ? R(e) : e;
			return L(v, n), _ = !0, c !== void 0 && (c = n), e;
		}
		return Pn && _ || y.f & 16384 ? v.v : K(v);
	});
}
//#endregion
//#region node_modules/svelte/src/internal/disclose-version.js
typeof window < "u" && ((window.__svelte ??= {}).v ??= /* @__PURE__ */ new Set()).add("5");
//#endregion
//#region src/lib/markup.js
function ni(e, t, n, r) {
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
function ri(e, t) {
	return e * (1 + t / 100);
}
function ii(e, t, n) {
	return e * ri(t, n);
}
function ai(e, t) {
	let n = {
		materials: 0,
		labor: 0,
		equipment: 0,
		subs: 0,
		other: 0
	}, r = 0, i = 0, a = (a) => {
		let o = ni(a.category_type, t, e.markup_overrides, e.markup_enabled), s = a.quantity * a.unit_price, c = ii(a.quantity, a.unit_price, o);
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
function oi(e, t) {
	let n = {
		materials: 0,
		labor: 0,
		equipment: 0,
		subs: 0,
		other: 0
	}, r = 0, i = 0;
	for (let a of e.subcategories) {
		let e = ai(a, t);
		r += e.base, i += e.withMarkup;
		for (let t of Object.keys(n)) n[t] += e.byType[t];
	}
	return {
		base: r,
		withMarkup: i,
		byType: n
	};
}
function si(e) {
	let t = {
		materials: 0,
		labor: 0,
		equipment: 0,
		subs: 0,
		other: 0
	}, n = 0, r = 0;
	for (let i of e.sections) {
		let a = oi(i, e.globals);
		n += a.base, r += a.withMarkup;
		for (let e of Object.keys(t)) t[e] += a.byType[e];
	}
	return {
		base: n,
		withMarkup: r,
		byType: t
	};
}
function ci(e) {
	return new Intl.NumberFormat("en-US", {
		style: "currency",
		currency: "USD",
		minimumFractionDigits: 2
	}).format(e);
}
function li(e) {
	return `${e}%`;
}
//#endregion
//#region src/lib/LineItemRow.svelte
var ui = /* @__PURE__ */ J("<option> </option>"), di = /* @__PURE__ */ J("<tr class=\"border-b border-white/[0.04]\"><td colspan=\"9\" class=\"px-2 py-1\"><div class=\"flex items-center gap-2 pl-6\"><svg class=\"w-3 h-3 text-white/15 shrink-0\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M3 10h10a3 3 0 013 3v1\"></path></svg> <input type=\"text\" placeholder=\"Add a note...\" class=\"w-full px-1 py-0.5 text-xs bg-transparent border-0 border-b border-white/[0.04] text-[var(--color-concrete)]\n						focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/15 font-[var(--font-body)]\"/></div></td></tr>"), fi = /* @__PURE__ */ J("<tr class=\"border-b border-white/[0.04] hover:bg-white/[0.02] text-sm group\"><td class=\"px-1 py-1 w-24\"><select class=\"w-full text-xs px-1 py-1 border-0 bg-transparent font-medium cursor-pointer\n				focus:ring-1 focus:ring-[var(--color-sunburst)] text-[var(--color-concrete)] font-[var(--font-ui)]\" style=\"color-scheme: dark;\"></select></td><td class=\"px-1 py-1\"><input type=\"text\" class=\"w-full px-1 py-0.5 text-[var(--color-white)] bg-transparent border-0\n				focus:ring-1 focus:ring-[var(--color-sunburst)] font-[var(--font-body)]\" placeholder=\"Item name\"/></td><td class=\"px-1 py-1 w-20\"><input type=\"number\" step=\"any\" min=\"0\" class=\"w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 text-[var(--color-white)]\n				focus:ring-1 focus:ring-[var(--color-sunburst)]\"/></td><td class=\"px-1 py-1 w-16\"><input type=\"text\" class=\"w-full text-center text-white/30 px-1 py-0.5 bg-transparent border-0\n				focus:ring-1 focus:ring-[var(--color-sunburst)] text-sm font-[var(--font-body)]\" placeholder=\"ea\"/></td><td class=\"px-1 py-1 w-24\"><input type=\"text\" inputmode=\"decimal\" class=\"w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 text-[var(--color-white)]\n				focus:ring-1 focus:ring-[var(--color-sunburst)]\"/></td><td class=\"px-2 py-1.5 text-right font-mono text-white/25 w-16 text-xs\"> </td><td class=\"px-2 py-1.5 text-right font-mono text-white/40 w-24 text-xs\"> </td><td class=\"px-2 py-1.5 text-right font-mono font-medium text-[var(--color-white)] w-24\"> </td><td class=\"px-1 py-1 w-16\"><div class=\"flex items-center gap-0.5\"><button title=\"Toggle note\"><svg class=\"w-3.5 h-3.5\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z\"></path></svg></button> <button class=\"opacity-0 group-hover:opacity-100 text-white/20 hover:text-red-400 transition-opacity p-0.5\" title=\"Delete item\"><svg class=\"w-3.5 h-3.5\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></button></div></td></tr> <!>", 1);
function pi(e, t) {
	Fe(t, !0);
	let n = $(t, "item", 7), r = [
		"materials",
		"labor",
		"equipment",
		"subs",
		"other"
	], i = /* @__PURE__ */ F(() => ni(n().category_type, t.globals, t.markupOverrides, t.markupEnabled)), a = /* @__PURE__ */ F(() => ri(n().unit_price, K(i))), o = /* @__PURE__ */ F(() => ii(n().quantity, n().unit_price, K(i))), s = /* @__PURE__ */ I(!!n().description);
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
		L(s, !K(s));
	}
	function g() {
		t.ondelete?.(n().id);
	}
	var _ = fi(), v = Qt(_), y = z(v), b = z(y);
	kr(b, 21, () => r, Tr, (e, t) => {
		var n = ui(), r = z(n, !0);
		A(n);
		var i = {};
		V(() => {
			X(r, K(t)), i !== (i = K(t)) && (n.value = (n.__value = K(t)) ?? "");
		}), Y(e, n);
	}), A(b);
	var x;
	zr(b), A(y);
	var S = B(y), ee = z(S);
	Q(ee), A(S);
	var te = B(S), C = z(te);
	Q(C), A(te);
	var ne = B(te), re = z(ne);
	Q(re), A(ne);
	var w = B(ne), T = z(w);
	Q(T), A(w);
	var ie = B(w), ae = z(ie);
	A(ie);
	var E = B(ie), oe = z(E, !0);
	A(E);
	var se = B(E), ce = z(se, !0);
	A(se);
	var le = B(se), ue = z(le), de = z(ue), fe = B(de, 2);
	A(ue), A(le), A(v);
	var pe = B(v, 2), me = (e) => {
		var t = di(), r = z(t), i = z(r), a = B(z(i), 2);
		Q(a), A(i), A(r), A(t), V(() => Gr(a, n().description ?? "")), q("input", a, m), Y(e, t);
	};
	Z(pe, (e) => {
		K(s) && e(me);
	}), V((e, t, r) => {
		x !== (x = n().category_type) && (b.value = (b.__value = n().category_type) ?? "", Rr(b, n().category_type)), Gr(ee, n().item_name), Gr(C, n().quantity), Gr(re, n().unit), Gr(T, e), X(ae, `${K(i) ?? ""}%`), X(oe, t), X(ce, r), Lr(de, 1, `p-0.5 transition-opacity
					${K(s) || n().description ? "opacity-100 text-[var(--color-sunburst)]" : "opacity-0 group-hover:opacity-100 text-white/20 hover:text-[var(--color-concrete)]"}`);
	}, [
		() => n().unit_price.toFixed(2),
		() => ci(K(a)),
		() => ci(K(o))
	]), q("change", b, p), q("input", ee, d), q("input", C, l), q("input", re, f), q("input", T, u), q("click", de, h), q("click", fe, g), Y(e, _), Ie();
}
fr([
	"change",
	"input",
	"click"
]);
//#endregion
//#region src/lib/Autocomplete.svelte
var mi = /* @__PURE__ */ J("<span class=\"text-xs text-[var(--color-sage)] font-[var(--font-ui)]\">Sub</span>"), hi = /* @__PURE__ */ J("<div class=\"flex items-center gap-2 text-xs\"><span class=\"font-mono text-[var(--color-sunburst)]\"> </span> <span class=\"text-white/30\"> </span></div>"), gi = /* @__PURE__ */ J("<button><div><span class=\"text-[var(--color-white)] font-[var(--font-body)]\"> </span> <span class=\"text-xs text-white/30 ml-2 font-[var(--font-body)]\"> </span></div> <!></button>"), _i = /* @__PURE__ */ J("<div class=\"absolute top-full left-0 right-0 mt-1 border border-white/[0.08] shadow-lg z-50 max-h-64 overflow-y-auto\" style=\"background: var(--color-granite);\"></div>"), vi = /* @__PURE__ */ J("<div class=\"absolute top-full left-0 right-0 mt-1 border border-white/[0.08] shadow-lg z-50\" style=\"background: var(--color-granite);\"><button class=\"w-full text-left px-3 py-2 text-sm text-[var(--color-muted-text)] hover:bg-white/[0.04] font-[var(--font-body)] transition-colors\">Add custom item: <span class=\"font-medium text-[var(--color-white)]\"> </span></button></div>"), yi = /* @__PURE__ */ J("<div class=\"relative\"><input type=\"text\" placeholder=\"Search items or type a name...\" class=\"w-full px-3 py-2 text-sm bg-white/[0.04] border border-white/[0.08] text-[var(--color-white)]\n			focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] outline-none placeholder-white/20 font-[var(--font-body)]\"/> <!></div>");
function bi(e, t) {
	Fe(t, !0);
	let n = $(t, "materialsDb", 19, () => []), r = $(t, "ratesDb", 19, () => []), i = $(t, "subcontractorsDb", 19, () => []), a = $(t, "categoryType", 3, "materials"), o = /* @__PURE__ */ I(""), s = /* @__PURE__ */ I(0), c, l = /* @__PURE__ */ F(() => {
		if (K(o).length < 1) return [];
		let e = K(o).toLowerCase(), t;
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
		e.key === "ArrowDown" ? (e.preventDefault(), L(s, Math.min(K(s) + 1, K(l).length - 1), !0)) : e.key === "ArrowUp" ? (e.preventDefault(), L(s, Math.max(K(s) - 1, 0), !0)) : e.key === "Enter" ? (e.preventDefault(), K(l).length > 0 && K(s) < K(l).length ? d(K(l)[K(s)]) : K(o).trim() && f()) : e.key === "Escape" && (e.preventDefault(), t.oncancel?.());
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
			item_name: K(o).trim(),
			unit: "ea",
			unit_price: 0,
			category_type: a(),
			is_custom: !0,
			material_id: null,
			subcontractor_id: null
		});
	}
	mn(() => {
		L(s, 0);
	}), mn(() => {
		c?.focus();
	});
	var p = yi(), m = z(p);
	Q(m), ti(m, (e) => c = e, () => c);
	var h = B(m, 2), g = (e) => {
		var t = _i();
		kr(t, 21, () => K(l), Tr, (e, t, n) => {
			var r = gi(), i = z(r), a = z(i), o = z(a, !0);
			A(a);
			var c = B(a, 2), l = z(c, !0);
			A(c), A(i);
			var u = B(i, 2), f = (e) => {
				Y(e, mi());
			}, p = (e) => {
				var n = hi(), r = z(n), i = z(r);
				A(r);
				var a = B(r, 2), o = z(a);
				A(a), A(n), V((e) => {
					X(i, `$${e ?? ""}`), X(o, `/ ${K(t).unit ?? ""}`);
				}, [() => K(t).price.toFixed(2)]), Y(e, n);
			};
			Z(u, (e) => {
				K(t).source === "subcontractor" ? e(f) : e(p, -1);
			}), A(r), V(() => {
				Lr(r, 1, `w-full text-left px-3 py-2 text-sm flex items-center justify-between
						hover:bg-white/[0.04] ${n === K(s) ? "bg-white/[0.06]" : ""} transition-colors`), X(o, K(t).name), X(l, K(t).category);
			}), q("click", r, () => d(K(t))), dr("mouseenter", r, () => L(s, n, !0)), Y(e, r);
		}), A(t), Y(e, t);
	}, _ = (e) => {
		var t = vi(), n = z(t), r = B(z(n)), i = z(r);
		A(r), A(n), A(t), V(() => X(i, `"${K(o) ?? ""}"`)), q("click", n, f), Y(e, t);
	};
	Z(h, (e) => {
		K(l).length > 0 ? e(g) : K(o).length > 0 && e(_, 1);
	}), A(p), q("keydown", m, u), Zr(m, () => K(o), (e) => L(o, e)), Y(e, p), Ie();
}
fr(["keydown", "click"]);
//#endregion
//#region node_modules/nanoid/url-alphabet/index.js
var xi = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict", Si = (e = 21) => {
	let t = "", n = crypto.getRandomValues(new Uint8Array(e |= 0));
	for (; e--;) t += xi[n[e] & 63];
	return t;
}, Ci = /* @__PURE__ */ J("<input type=\"text\" class=\"text-xs font-semibold text-[var(--color-concrete)] uppercase tracking-wide px-1 py-0.5 border border-white/[0.08] bg-white/[0.04] focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] font-[var(--font-ui)]\"/>"), wi = /* @__PURE__ */ J("<span class=\"text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider cursor-pointer font-[var(--font-ui)]\" role=\"button\" tabindex=\"0\"> </span> <button class=\"opacity-0 group-hover/cg:opacity-100 text-white/15 hover:text-[var(--color-sunburst)] transition-opacity p-0.5\" title=\"Rename\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z\"></path></svg></button>", 1), Ti = /* @__PURE__ */ J("<table class=\"w-full\"><tbody></tbody></table>"), Ei = /* @__PURE__ */ J("<div class=\"mt-1 mb-1\"><!></div>"), Di = /* @__PURE__ */ J("<button class=\"mt-1 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Item</button>"), Oi = /* @__PURE__ */ J("<div class=\"ml-5 mt-3 border-l-2 border-white/[0.06] pl-4\"><div class=\"flex items-center gap-2 mb-1 group/cg\"><!> <span class=\"text-xs text-white/25 font-[var(--font-body)]\"> </span> <button class=\"opacity-0 group-hover/cg:opacity-100 text-white/20 hover:text-red-400 transition-opacity p-0.5\" title=\"Delete group\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></button></div> <!> <!></div>");
function ki(e, t) {
	Fe(t, !0);
	let n = $(t, "group", 7), r = /* @__PURE__ */ I(!1), i = /* @__PURE__ */ I(!1), a = /* @__PURE__ */ I(R(n().name));
	function o(e) {
		t.onsnapshot?.(), n().line_items.push({
			id: Si(),
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
			component_group_id: n().id
		}), L(r, !1), t.onchange?.();
	}
	function s(e) {
		t.onsnapshot?.();
		let r = n().line_items.findIndex((t) => t.id === e);
		r !== -1 && (n().line_items.splice(r, 1), t.onchange?.());
	}
	function c() {
		L(a, n().name, !0), L(i, !0);
	}
	function l() {
		let e = K(a).trim();
		e && e !== n().name && (t.onsnapshot?.(), n().name = e, t.onchange?.()), L(i, !1);
	}
	var u = Oi(), d = z(u), f = z(d), p = (e) => {
		var t = Ci();
		Q(t), rn(t, !0), q("keydown", t, (e) => {
			e.key === "Enter" && l(), e.key === "Escape" && L(i, !1);
		}), dr("blur", t, l), Zr(t, () => K(a), (e) => L(a, e)), Y(e, t);
	}, m = (e) => {
		var t = wi(), r = Qt(t), i = z(r, !0);
		A(r);
		var a = B(r, 2);
		V(() => X(i, n().name)), q("dblclick", r, c), q("click", a, c), Y(e, t);
	};
	Z(f, (e) => {
		K(i) ? e(p) : e(m, -1);
	});
	var h = B(f, 2), g = z(h);
	A(h);
	var _ = B(h, 2);
	A(d);
	var v = B(d, 2), y = (e) => {
		var r = Ti(), i = z(r);
		kr(i, 21, () => n().line_items, (e) => e.id, (e, n) => {
			pi(e, {
				get item() {
					return K(n);
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
				ondelete: s
			});
		}), A(i), A(r), Y(e, r);
	};
	Z(v, (e) => {
		n().line_items.length > 0 && e(y);
	});
	var b = B(v, 2), x = (e) => {
		var n = Ei();
		bi(z(n), {
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
			onselect: o,
			oncancel: () => L(r, !1)
		}), A(n), Y(e, n);
	}, S = (e) => {
		var t = Di();
		q("click", t, () => L(r, !0)), Y(e, t);
	};
	Z(b, (e) => {
		K(r) ? e(x) : e(S, -1);
	}), A(u), V(() => X(g, `(${n().line_items.length ?? ""})`)), q("click", _, () => t.ondelete?.(n().id)), Y(e, u), Ie();
}
fr([
	"keydown",
	"dblclick",
	"click"
]);
//#endregion
//#region src/lib/SubcategoryBlock.svelte
var Ai = /* @__PURE__ */ J("<input type=\"text\" class=\"px-2 py-0.5 border border-white/[0.08] text-sm font-medium text-[var(--color-white)] bg-white/[0.04] focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] font-[var(--font-ui)]\"/>"), ji = /* @__PURE__ */ J("<span class=\"font-medium text-[var(--color-concrete)] text-sm font-[var(--font-ui)]\"> </span> <button class=\"opacity-0 group-hover/subcat:opacity-100 text-white/15 hover:text-[var(--color-sunburst)] transition-opacity p-0.5 shrink-0\" title=\"Rename\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z\"></path></svg></button>", 1), Mi = /* @__PURE__ */ J("<span class=\"text-xs px-1.5 py-0.5 bg-[var(--color-sunburst)]/10 text-[var(--color-sunburst)] font-medium font-[var(--font-ui)]\">overrides</span>"), Ni = /* @__PURE__ */ J("<span class=\"text-xs px-1.5 py-0.5 bg-[var(--color-sage)]/10 text-[var(--color-sage)] font-medium font-[var(--font-ui)]\"> </span>"), Pi = /* @__PURE__ */ J("<div class=\"text-center\"><span class=\"block text-xs font-medium text-[var(--color-concrete)] mb-1 font-[var(--font-ui)]\"> </span> <div class=\"flex items-center justify-center gap-1 mb-1\"><input type=\"checkbox\" class=\"w-3 h-3 border-white/[0.1] bg-white/[0.04] text-[var(--color-sunburst)] focus:ring-[var(--color-sunburst)] rounded-sm\"/> <span class=\"text-xs text-white/30 font-[var(--font-body)]\"> </span></div> <input type=\"number\" step=\"1\" min=\"0\"/> <div class=\"text-xs text-white/30 mt-0.5 font-mono\"> </div></div>"), Fi = /* @__PURE__ */ J("<div class=\"bg-white/[0.02] border border-white/[0.06] p-3 mb-3\"><div class=\"flex items-center justify-between mb-2\"><span class=\"text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider font-[var(--font-ui)]\">Markup Overrides</span> <button class=\"text-xs text-white/30 hover:text-[var(--color-white)] font-[var(--font-ui)]\">Close</button></div> <div class=\"grid grid-cols-5 gap-3\"></div> <div class=\"mt-3 pt-2 border-t border-white/[0.06] flex items-center gap-2\"><span class=\"text-xs font-medium text-[var(--color-concrete)] font-[var(--font-ui)]\">Lump Sum:</span> <span class=\"text-xs text-white/30\">$</span> <input type=\"number\" step=\"0.01\" min=\"0\" class=\"w-28 text-right text-xs font-mono px-2 py-1 border border-white/[0.08] bg-white/[0.04] text-[var(--color-white)]\n								focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)]\"/> <span class=\"text-xs text-white/30 font-[var(--font-body)]\">added post-markup</span></div></div>"), Ii = /* @__PURE__ */ J("<span> </span>"), Li = /* @__PURE__ */ J("<div class=\"flex items-center gap-3 py-1.5 px-3 bg-[var(--color-sunburst)]/5 text-xs mb-2 cursor-pointer hover:bg-[var(--color-sunburst)]/8 transition-colors border-l-2 border-[var(--color-sunburst)]/30\" role=\"button\" tabindex=\"0\"><span class=\"font-medium text-[var(--color-sunburst)] font-[var(--font-ui)]\">Markup:</span> <!></div>"), Ri = /* @__PURE__ */ J("<table class=\"w-full\"><thead><tr class=\"text-xs text-white/30 uppercase tracking-wider border-b border-white/[0.06] font-[var(--font-ui)]\"><th class=\"px-1 py-1.5 text-left w-24\">Type</th><th class=\"px-1 py-1.5 text-left\">Name</th><th class=\"px-1 py-1.5 text-right w-20\">Qty</th><th class=\"px-1 py-1.5 text-center w-16\">Unit</th><th class=\"px-1 py-1.5 text-right w-24\">Price</th><th class=\"px-2 py-1.5 text-right w-16\">Markup</th><th class=\"px-2 py-1.5 text-right w-24\">w/ Markup</th><th class=\"px-2 py-1.5 text-right w-24\">Total</th><th class=\"w-16\"></th></tr></thead><tbody></tbody></table>"), zi = /* @__PURE__ */ J("<div class=\"flex items-center gap-2 ml-5 mt-2\"><input type=\"text\" placeholder=\"Group name\" class=\"flex-1 px-0 py-1 bg-transparent border-0 border-b border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/20\"/> <button class=\"px-2 py-1 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-xs font-[var(--font-ui)] font-bold uppercase tracking-wide hover:brightness-110\">Add</button> <button class=\"px-2 py-1 text-[var(--color-muted-text)] text-xs hover:text-[var(--color-white)] font-[var(--font-ui)]\">Cancel</button></div>"), Bi = /* @__PURE__ */ J("<button class=\"ml-5 mt-2 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Group</button>"), Vi = /* @__PURE__ */ J("<div class=\"mt-2\"><!></div>"), Hi = /* @__PURE__ */ J("<button class=\"mt-2 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors\"><svg class=\"w-3.5 h-3.5\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Item</button>"), Ui = /* @__PURE__ */ J("<span class=\"text-[var(--color-sage)]\"> </span>"), Wi = /* @__PURE__ */ J("<div class=\"px-5 pb-3\"><!> <!> <!> <!> <!> <div class=\"flex justify-between items-center mt-3 pt-2 border-t border-white/[0.06] text-sm\"><span class=\"text-white/30 font-[var(--font-body)]\"> <!></span> <span class=\"font-mono font-semibold text-[var(--color-white)]\"> </span></div></div>"), Gi = /* @__PURE__ */ J("<div class=\"border-t border-white/[0.06]\"><div class=\"flex items-center justify-between px-5 py-2 hover:bg-white/[0.02] transition-colors group/subcat\"><button class=\"flex items-center gap-2 text-left flex-1 min-w-0\"><svg fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> <!> <span class=\"text-xs text-white/25 font-[var(--font-body)]\"> </span> <!> <!></button> <div class=\"flex items-center gap-2\"><button class=\"text-xs px-1.5 py-0.5 bg-white/[0.04] text-[var(--color-muted-text)] hover:bg-white/[0.08] hover:text-[var(--color-concrete)] transition-colors\" title=\"Configure markup\"><svg class=\"w-3 h-3 inline\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4\"></path></svg></button> <span class=\"font-mono text-sm font-semibold text-[var(--color-white)]\"> </span> <button class=\"opacity-0 group-hover/subcat:opacity-100 text-white/20 hover:text-red-400 transition-opacity p-0.5\" title=\"Delete subcategory\"><svg class=\"w-3.5 h-3.5\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></button></div></div> <!></div>");
function Ki(e, t) {
	Fe(t, !0);
	let n = $(t, "subcat", 7), r = $(t, "collapsed", 3, !1), i = $(t, "materialsDb", 19, () => []), a = $(t, "ratesDb", 19, () => []), o = $(t, "subcontractorsDb", 19, () => []), s = /* @__PURE__ */ I(R(r())), c = /* @__PURE__ */ I(!1), l = /* @__PURE__ */ I(!1), u = /* @__PURE__ */ I(R(n().name)), d = /* @__PURE__ */ I(!1), f = /* @__PURE__ */ I(""), p = /* @__PURE__ */ F(() => ai(n(), t.globals)), m = /* @__PURE__ */ F(() => [
		"materials",
		"labor",
		"equipment",
		"subs",
		"other"
	].map((e) => ({
		type: e,
		value: ni(e, t.globals, n().markup_overrides, n().markup_enabled),
		isOverride: n().markup_overrides[e] != null,
		isDisabled: !n().markup_enabled[e]
	}))), h = /* @__PURE__ */ I(!1), g = /* @__PURE__ */ F(() => K(m).some((e) => e.isOverride || e.isDisabled)), _ = /* @__PURE__ */ F(() => n().line_items.length + n().component_groups.reduce((e, t) => e + t.line_items.length, 0));
	function v(e) {
		t.onsnapshot?.(), n().line_items.push({
			id: Si(),
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
			component_group_id: null
		}), L(c, !1), t.onchange?.();
	}
	function y(e) {
		t.onsnapshot?.();
		let r = n().line_items.findIndex((t) => t.id === e);
		r !== -1 && (n().line_items.splice(r, 1), t.onchange?.());
	}
	function b() {
		L(u, n().name, !0), L(l, !0);
	}
	function x() {
		let e = K(u).trim();
		e && e !== n().name && (t.onsnapshot?.(), n().name = e, t.onchange?.()), L(l, !1);
	}
	function S() {
		t.onsnapshot?.();
		let e = K(f).trim() || "New Group";
		n().component_groups.push({
			id: Si(),
			name: e,
			sort_order: n().component_groups.length,
			line_items: []
		}), L(f, ""), L(d, !1), t.onchange?.();
	}
	function ee(e) {
		t.onsnapshot?.();
		let r = n().component_groups.findIndex((t) => t.id === e);
		r !== -1 && (n().component_groups.splice(r, 1), t.onchange?.());
	}
	function te(e, r) {
		let i = r.target.value.trim();
		if (i === "") n().markup_overrides[e] = null;
		else {
			let t = parseFloat(i);
			!isNaN(t) && t >= 0 && (n().markup_overrides[e] = t);
		}
		t.onchange?.();
	}
	function C(e) {
		n().markup_enabled[e] = !n().markup_enabled[e], t.onchange?.();
	}
	function ne(e) {
		let r = parseFloat(e.target.value);
		!isNaN(r) && r >= 0 && (n().lump_sum = r, t.onchange?.());
	}
	let re = {
		materials: "Mat",
		labor: "Lab",
		equipment: "Equip",
		subs: "Subs",
		other: "Other"
	};
	var w = Gi(), T = z(w), ie = z(T), ae = z(ie), E = B(ae, 2), oe = (e) => {
		var t = Ai();
		Q(t), rn(t, !0), q("click", t, (e) => e.stopPropagation()), q("keydown", t, (e) => {
			e.key === "Enter" && x(), e.key === "Escape" && L(l, !1);
		}), dr("blur", t, x), Zr(t, () => K(u), (e) => L(u, e)), Y(e, t);
	}, se = (e) => {
		var t = ji(), r = Qt(t), i = z(r, !0);
		A(r);
		var a = B(r, 2);
		V(() => X(i, n().name)), q("click", a, (e) => {
			e.stopPropagation(), b();
		}), Y(e, t);
	};
	Z(E, (e) => {
		K(l) ? e(oe) : e(se, -1);
	});
	var ce = B(E, 2), le = z(ce);
	A(ce);
	var ue = B(ce, 2), de = (e) => {
		Y(e, Mi());
	};
	Z(ue, (e) => {
		K(g) && e(de);
	});
	var fe = B(ue, 2), pe = (e) => {
		var t = Ni(), r = z(t);
		A(t), V((e) => X(r, `+${e ?? ""} lump sum`), [() => ci(n().lump_sum)]), Y(e, t);
	};
	Z(fe, (e) => {
		n().lump_sum > 0 && e(pe);
	}), A(ie);
	var me = B(ie, 2), he = z(me), ge = B(he, 2), _e = z(ge, !0);
	A(ge);
	var ve = B(ge, 2);
	A(me), A(T);
	var D = B(T, 2), ye = (e) => {
		var r = Wi(), s = z(r), l = (e) => {
			var r = Fi(), i = z(r), a = B(z(i), 2);
			A(i);
			var o = B(i, 2);
			kr(o, 20, () => [
				"materials",
				"labor",
				"equipment",
				"subs",
				"other"
			], Tr, (e, r) => {
				let i = /* @__PURE__ */ F(() => K(m).find((e) => e.type === r));
				var a = Pi(), o = z(a), s = z(o, !0);
				A(o);
				var c = B(o, 2), l = z(c);
				Q(l);
				var u = B(l, 2), d = z(u, !0);
				A(u), A(c);
				var f = B(c, 2);
				Q(f);
				var p = B(f, 2), h = z(p);
				A(p), A(a), V((e) => {
					X(s, re[r]), Kr(l, n().markup_enabled[r]), X(d, n().markup_enabled[r] ? "On" : "Off"), Gr(f, n().markup_overrides[r] ?? ""), qr(f, "placeholder", `${t.globals[`${r}_markup`]}%`), f.disabled = !n().markup_enabled[r], Lr(f, 1, `w-full text-center text-xs font-mono px-1 py-1 border border-white/[0.08]
										focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)]
										${n().markup_enabled[r] ? "bg-white/[0.04] text-[var(--color-white)]" : "bg-white/[0.01] text-white/15"}
										${K(i)?.isOverride ? "border-[var(--color-sunburst)]/40 bg-[var(--color-sunburst)]/5" : ""}`), X(h, `eff: ${e ?? ""}`);
				}, [() => li(K(i)?.value ?? 0)]), q("change", l, () => C(r)), q("input", f, (e) => te(r, e)), Y(e, a);
			}), A(o);
			var s = B(o, 2), c = B(z(s), 4);
			Q(c), Ee(2), A(s), A(r), V(() => Gr(c, n().lump_sum)), q("click", a, () => L(h, !1)), q("input", c, ne), Y(e, r);
		}, u = (e) => {
			var t = Li();
			kr(B(z(t), 2), 17, () => K(m), Tr, (e, t) => {
				var n = Ii(), r = z(n);
				A(n), V((e) => {
					Lr(n, 1, `${K(t).isDisabled ? "text-white/15 line-through" : K(t).isOverride ? "text-[var(--color-sunburst)]" : "text-white/30"} font-[var(--font-body)]`), X(r, `${K(t).type ?? ""} ${e ?? ""}`);
				}, [() => li(K(t).value)]), Y(e, n);
			}), A(t), q("click", t, () => L(h, !0)), q("keydown", t, (e) => {
				(e.key === "Enter" || e.key === " ") && L(h, !0);
			}), Y(e, t);
		};
		Z(s, (e) => {
			K(h) ? e(l) : K(g) && e(u, 1);
		});
		var b = B(s, 2), x = (e) => {
			var r = Ri(), i = B(z(r));
			kr(i, 21, () => n().line_items, (e) => e.id, (e, r) => {
				pi(e, {
					get item() {
						return K(r);
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
					ondelete: y
				});
			}), A(i), A(r), Y(e, r);
		};
		Z(b, (e) => {
			(K(_) > 0 || n().line_items.length > 0) && e(x);
		});
		var w = B(b, 2);
		kr(w, 17, () => n().component_groups, (e) => e.id, (e, r) => {
			ki(e, {
				get group() {
					return K(r);
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
				ondelete: ee,
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
		var T = B(w, 2), ie = (e) => {
			var t = zi(), n = z(t);
			Q(n);
			var r = B(n, 2), i = B(r, 2);
			A(t), q("keydown", n, (e) => {
				e.key === "Enter" && S(), e.key === "Escape" && L(d, !1);
			}), Zr(n, () => K(f), (e) => L(f, e)), q("click", r, S), q("click", i, () => L(d, !1)), Y(e, t);
		}, ae = (e) => {
			var t = Bi();
			q("click", t, () => L(d, !0)), Y(e, t);
		};
		Z(T, (e) => {
			K(d) ? e(ie) : e(ae, -1);
		});
		var E = B(T, 2), oe = (e) => {
			var t = Vi();
			bi(z(t), {
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
				onselect: v,
				oncancel: () => L(c, !1)
			}), A(t), Y(e, t);
		}, se = (e) => {
			var t = Hi();
			q("click", t, () => L(c, !0)), Y(e, t);
		};
		Z(E, (e) => {
			K(c) ? e(oe) : e(se, -1);
		});
		var ce = B(E, 2), le = z(ce), ue = z(le), de = B(ue), fe = (e) => {
			var t = Ui(), r = z(t);
			A(t), V((e) => X(r, `+ ${e ?? ""} lump sum`), [() => ci(n().lump_sum)]), Y(e, t);
		};
		Z(de, (e) => {
			n().lump_sum > 0 && e(fe);
		}), A(le);
		var pe = B(le, 2), me = z(pe, !0);
		A(pe), A(ce), A(r), V((e, t) => {
			X(ue, `Subtotal: ${e ?? ""} `), X(me, t);
		}, [() => ci(K(p).base), () => ci(K(p).withMarkup)]), Y(e, r);
	};
	Z(D, (e) => {
		K(s) || e(ye);
	}), A(w), V((e) => {
		Lr(ae, 0, `w-3 h-3 text-white/25 transition-transform shrink-0 ${K(s) ? "" : "rotate-90"}`), X(le, `${K(_) ?? ""} item${K(_) === 1 ? "" : "s"}`), X(_e, e);
	}, [() => ci(K(p).withMarkup)]), q("click", ie, () => L(s, !K(s))), q("click", he, () => L(h, !K(h))), q("click", ve, () => t.ondelete?.(n().id)), Y(e, w), Ie();
}
fr([
	"click",
	"keydown",
	"change",
	"input"
]);
//#endregion
//#region src/lib/SectionBlock.svelte
var qi = /* @__PURE__ */ J("<input type=\"text\" class=\"bg-white/[0.06] text-[var(--color-white)] px-2 py-0.5 text-sm font-semibold border border-white/[0.08] focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] font-[var(--font-ui)] uppercase tracking-wide\"/>"), Ji = /* @__PURE__ */ J("<span class=\"font-semibold font-[var(--font-ui)] uppercase tracking-wide text-[var(--color-white)]\"> </span> <button class=\"opacity-0 group-hover/sec:opacity-100 text-white/20 hover:text-[var(--color-sunburst)] transition-opacity p-0.5 shrink-0\" title=\"Rename section\"><svg class=\"w-3.5 h-3.5\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z\"></path></svg></button>", 1), Yi = /* @__PURE__ */ J("<div class=\"px-5 py-10 text-center text-[var(--color-muted-text)] text-sm font-[var(--font-body)]\">No subcategories yet</div>"), Xi = /* @__PURE__ */ J("<div class=\"flex items-center gap-2 mt-2\"><input type=\"text\" placeholder=\"Subcategory name\" class=\"flex-1 px-0 py-1.5 bg-transparent border-0 border-b border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/20\"/> <button class=\"px-3 py-1.5 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-xs font-[var(--font-ui)] font-bold uppercase tracking-wide hover:brightness-110\">Add</button> <button class=\"px-2 py-1.5 text-[var(--color-muted-text)] text-xs hover:text-[var(--color-white)] font-[var(--font-ui)]\">Cancel</button></div>"), Zi = /* @__PURE__ */ J("<button class=\"mt-2 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] uppercase tracking-wide transition-colors\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Subcategory</button>"), Qi = /* @__PURE__ */ J("<!> <div class=\"px-5 pb-4\"><!></div>", 1), $i = /* @__PURE__ */ J("<div class=\"mb-6 border border-white/[0.06] overflow-hidden\" style=\"background: var(--color-ink);\"><div class=\"flex items-center justify-between px-5 py-3 group/sec\" style=\"background: var(--color-granite);\"><button class=\"flex items-center gap-3 hover:bg-white/[0.04] -ml-2 px-2 py-1 transition-colors text-left flex-1 min-w-0\"><svg fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> <!> <span class=\"text-xs text-white/30 font-[var(--font-body)]\"> </span></button> <div class=\"flex items-center gap-3\"><span class=\"font-mono font-semibold text-[var(--color-white)]\"> </span> <button class=\"text-white/20 hover:text-red-400 transition-colors p-1\" title=\"Delete section\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16\"></path></svg></button></div></div> <!></div>");
function ea(e, t) {
	Fe(t, !0);
	let n = $(t, "section", 7), r = $(t, "collapsed", 3, !1), i = $(t, "materialsDb", 19, () => []), a = $(t, "ratesDb", 19, () => []), o = $(t, "subcontractorsDb", 19, () => []), s = /* @__PURE__ */ I(R(r())), c = /* @__PURE__ */ I(!1), l = /* @__PURE__ */ I(R(n().name)), u = /* @__PURE__ */ I(!1), d = /* @__PURE__ */ I(""), f = /* @__PURE__ */ F(() => oi(n(), t.globals)), p = /* @__PURE__ */ F(() => n().subcategories.reduce((e, t) => e + t.line_items.length + t.component_groups.reduce((e, t) => e + t.line_items.length, 0), 0));
	function m() {
		L(l, n().name, !0), L(c, !0);
	}
	function h() {
		let e = K(l).trim();
		e && e !== n().name && (t.onsnapshot?.(), n().name = e, t.onchange?.()), L(c, !1);
	}
	function g() {
		t.onsnapshot?.();
		let e = K(d).trim() || "New Subcategory";
		n().subcategories.push({
			id: Si(),
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
		}), L(d, ""), L(u, !1), t.onchange?.();
	}
	function _(e) {
		t.onsnapshot?.();
		let r = n().subcategories.findIndex((t) => t.id === e);
		r !== -1 && (n().subcategories.splice(r, 1), t.onchange?.());
	}
	var v = $i(), y = z(v), b = z(y), x = z(b), S = B(x, 2), ee = (e) => {
		var t = qi();
		Q(t), rn(t, !0), q("click", t, (e) => e.stopPropagation()), q("keydown", t, (e) => {
			e.key === "Enter" && h(), e.key === "Escape" && L(c, !1);
		}), dr("blur", t, h), Zr(t, () => K(l), (e) => L(l, e)), Y(e, t);
	}, te = (e) => {
		var t = Ji(), r = Qt(t), i = z(r, !0);
		A(r);
		var a = B(r, 2);
		V(() => X(i, n().name)), q("click", a, (e) => {
			e.stopPropagation(), m();
		}), Y(e, t);
	};
	Z(S, (e) => {
		K(c) ? e(ee) : e(te, -1);
	});
	var C = B(S, 2), ne = z(C);
	A(C), A(b);
	var re = B(b, 2), w = z(re), T = z(w, !0);
	A(w);
	var ie = B(w, 2);
	A(re), A(y);
	var ae = B(y, 2), E = (e) => {
		var r = Qi(), s = Qt(r), c = (e) => {
			Y(e, Yi());
		}, l = (e) => {
			var r = yr();
			kr(Qt(r), 17, () => n().subcategories, (e) => e.id, (e, n) => {
				Ki(e, {
					get subcat() {
						return K(n);
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
			}), Y(e, r);
		};
		Z(s, (e) => {
			n().subcategories.length === 0 && !K(u) ? e(c) : e(l, -1);
		});
		var f = B(s, 2), p = z(f), m = (e) => {
			var t = Xi(), n = z(t);
			Q(n);
			var r = B(n, 2), i = B(r, 2);
			A(t), q("keydown", n, (e) => {
				e.key === "Enter" && g(), e.key === "Escape" && L(u, !1);
			}), Zr(n, () => K(d), (e) => L(d, e)), q("click", r, g), q("click", i, () => L(u, !1)), Y(e, t);
		}, h = (e) => {
			var t = Zi();
			q("click", t, () => L(u, !0)), Y(e, t);
		};
		Z(p, (e) => {
			K(u) ? e(m) : e(h, -1);
		}), A(f), Y(e, r);
	};
	Z(ae, (e) => {
		K(s) || e(E);
	}), A(v), V((e) => {
		Lr(x, 0, `w-4 h-4 text-white/30 transition-transform shrink-0 ${K(s) ? "" : "rotate-90"}`), X(ne, `${n().subcategories.length ?? ""} subcategor${n().subcategories.length === 1 ? "y" : "ies"}
				· ${K(p) ?? ""} item${K(p) === 1 ? "" : "s"}`), X(T, e);
	}, [() => ci(K(f).withMarkup)]), q("click", b, () => L(s, !K(s))), q("click", ie, () => t.ondelete?.(n().id)), Y(e, v), Ie();
}
fr(["click", "keydown"]);
//#endregion
//#region src/lib/FooterSummary.svelte
var ta = /* @__PURE__ */ J("<div class=\"flex flex-col\"><span class=\"text-xs text-white/30 uppercase tracking-wider font-[var(--font-ui)]\"> </span> <span class=\"font-mono text-[var(--color-concrete)]\"> </span></div>"), na = /* @__PURE__ */ J("<span class=\"text-[var(--color-muted-text)] text-xs font-[var(--font-body)]\">No items yet</span>"), ra = /* @__PURE__ */ J("<div class=\"fixed bottom-0 right-0 border-t-2 z-20\" style=\"left: var(--sidebar-width, 220px); background: var(--color-granite); border-color: var(--color-sunburst);\"><div class=\"flex items-center justify-between px-5 py-3\"><div class=\"flex items-center gap-6 text-sm\"><div class=\"flex flex-col\"><span class=\"text-xs text-white/30 uppercase tracking-wider font-[var(--font-ui)]\">Base Cost</span> <span class=\"font-mono text-[var(--color-concrete)]\"> </span></div> <div class=\"w-px h-8 bg-white/[0.08]\"></div> <!> <!></div> <div class=\"flex flex-col items-end\"><span class=\"text-xs text-white/30 uppercase tracking-wider font-[var(--font-ui)]\">Total</span> <span class=\"font-mono text-lg font-bold text-[var(--color-sunburst)]\"> </span></div></div></div>");
function ia(e, t) {
	Fe(t, !0);
	let n = /* @__PURE__ */ F(() => si(t.estimate)), r = /* @__PURE__ */ F(() => Object.entries(K(n).byType).filter(([, e]) => e > 0).map(([e, t]) => ({
		type: e,
		value: t
	}))), i = {
		materials: "Materials",
		labor: "Labor",
		equipment: "Equipment",
		subs: "Subs",
		other: "Other"
	};
	var a = ra(), o = z(a), s = z(o), c = z(s), l = B(z(c), 2), u = z(l, !0);
	A(l), A(c);
	var d = B(c, 4);
	kr(d, 17, () => K(r), Tr, (e, t) => {
		let n = () => K(t).type, r = () => K(t).value;
		var a = ta(), o = z(a), s = z(o, !0);
		A(o);
		var c = B(o, 2), l = z(c, !0);
		A(c), A(a), V((e) => {
			X(s, i[n()]), X(l, e);
		}, [() => ci(r())]), Y(e, a);
	});
	var f = B(d, 2), p = (e) => {
		Y(e, na());
	};
	Z(f, (e) => {
		K(r).length === 0 && e(p);
	}), A(s);
	var m = B(s, 2), h = B(z(m), 2), g = z(h, !0);
	A(h), A(m), A(o), A(a), V((e, t) => {
		X(u, e), X(g, t);
	}, [() => ci(K(n).base), () => ci(K(n).withMarkup)]), Y(e, a), Ie();
}
//#endregion
//#region src/lib/SaveStatus.svelte
var aa = /* @__PURE__ */ J("<button class=\"ml-1 px-1.5 py-0.5 rounded bg-[var(--color-sunburst)]/20 hover:bg-[var(--color-sunburst)]/30 text-[var(--color-sunburst)] transition-colors\" title=\"Save now (Ctrl+S)\">Save</button>"), oa = /* @__PURE__ */ J("<button class=\"ml-1 px-1.5 py-0.5 rounded bg-red-500/20 hover:bg-red-500/30 text-red-300 transition-colors\" title=\"Retry save\">Retry</button>"), sa = /* @__PURE__ */ J("<div><span></span> <span> </span> <!> <!></div>");
function ca(e, t) {
	Fe(t, !0);
	let n = /* @__PURE__ */ I(R(Date.now()));
	mn(() => {
		let e = setInterval(() => {
			L(n, Date.now(), !0);
		}, 1e4);
		return () => clearInterval(e);
	});
	let r = /* @__PURE__ */ F(() => {
		switch (t.status) {
			case "clean":
				if (t.savedAt) {
					let e = Math.round((K(n) - t.savedAt.getTime()) / 1e3);
					return e < 5 ? "Saved just now" : e < 60 ? `Saved ${e}s ago` : `Saved ${Math.round(e / 60)}m ago`;
				}
				return "Up to date";
			case "dirty": return "Unsaved changes";
			case "saving": return "Saving...";
			case "error": return "Save failed";
			default: return "";
		}
	}), i = /* @__PURE__ */ F(() => {
		switch (t.status) {
			case "clean": return "text-[var(--color-sage)]";
			case "dirty": return "text-[var(--color-sunburst)]";
			case "saving": return "text-blue-400";
			case "error": return "text-red-400";
			default: return "text-white/40";
		}
	}), a = /* @__PURE__ */ F(() => {
		switch (t.status) {
			case "clean": return "bg-[var(--color-sage)]";
			case "dirty": return "bg-[var(--color-sunburst)]";
			case "saving": return "bg-blue-400 animate-pulse";
			case "error": return "bg-red-400 animate-pulse";
			default: return "bg-white/40";
		}
	});
	var o = sa(), s = z(o), c = B(s, 2), l = z(c, !0);
	A(c);
	var u = B(c, 2), d = (e) => {
		var n = aa();
		q("click", n, function(...e) {
			t.onsave?.apply(this, e);
		}), Y(e, n);
	};
	Z(u, (e) => {
		t.status === "dirty" && t.onsave && e(d);
	});
	var f = B(u, 2), p = (e) => {
		var n = oa();
		q("click", n, function(...e) {
			t.onsave?.apply(this, e);
		}), Y(e, n);
	};
	Z(f, (e) => {
		t.status === "error" && t.onsave && e(p);
	}), A(o), V(() => {
		Lr(o, 1, `flex items-center gap-1.5 text-xs font-[var(--font-ui)] ${K(i) ?? ""}`), Lr(s, 1, `w-2 h-2 rounded-full ${K(a) ?? ""}`), X(l, K(r));
	}), Y(e, o), Ie();
}
fr(["click"]);
//#endregion
//#region src/lib/autosave.svelte.js
function la(e, t = 2e3) {
	let n = /* @__PURE__ */ I("clean"), r = /* @__PURE__ */ I(null), i = null, a = null;
	function o(e) {
		a = e;
	}
	function s() {
		L(n, "dirty"), i && clearTimeout(i), i = setTimeout(() => c(), t);
	}
	async function c() {
		if (i &&= (clearTimeout(i), null), !a) return null;
		let o = a();
		if (!o) return null;
		L(n, "saving");
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
			return L(n, "clean"), L(r, /* @__PURE__ */ new Date(), !0), i;
		} catch (e) {
			return L(n, "error"), console.error("Auto-save failed:", e.message), i = setTimeout(() => c(), t * 2), null;
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
			return K(n);
		},
		get savedAt() {
			return K(r);
		}
	};
}
//#endregion
//#region src/lib/undo.svelte.js
var ua = 20;
function da() {
	let e = R([]), t = /* @__PURE__ */ F(() => e.length > 0);
	function n(t) {
		let n = JSON.parse(JSON.stringify({
			globals: t.globals,
			sections: t.sections
		}));
		e.push(n), e.length > ua && e.shift();
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
			return K(t);
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
var fa = /* @__PURE__ */ J("<div class=\"flex items-center justify-center h-64\"><div class=\"text-[var(--color-muted-text)] font-[var(--font-body)]\">Loading estimate...</div></div>"), pa = /* @__PURE__ */ J("<div class=\"flex items-center justify-center h-64\"><div class=\"text-red-400 font-[var(--font-body)]\"> </div></div>"), ma = /* @__PURE__ */ J("<label class=\"flex items-center gap-1 text-xs text-[var(--color-concrete)] font-[var(--font-ui)]\"><span> </span> <input type=\"number\" step=\"1\" min=\"0\" class=\"w-10 text-right bg-transparent border-0 border-b border-white/[0.08] p-0 font-mono text-xs text-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none\"/> <span class=\"text-white/30\">%</span></label>"), ha = /* @__PURE__ */ J("<div class=\"text-center py-16 text-[var(--color-muted-text)]\"><p class=\"text-lg font-[var(--font-ui)]\">No sections yet</p> <p class=\"text-sm mt-1 font-[var(--font-body)]\">Add a section to start building your estimate.</p> <button class=\"mt-4 px-4 py-2 bg-[var(--color-granite)] border border-white/[0.06] text-[var(--color-sunburst)] text-sm font-[var(--font-ui)] font-semibold uppercase tracking-wide hover:border-[var(--color-sunburst)] transition-colors\">+ Add First Section</button></div>"), ga = /* @__PURE__ */ J("<div class=\"flex items-center gap-2 mt-4\"><input type=\"text\" placeholder=\"Section name (e.g. Framing, Roofing)\" class=\"flex-1 px-0 py-2 bg-transparent border-0 border-b-2 border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/20\"/> <button class=\"px-4 py-2 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-xs font-[var(--font-ui)] font-bold uppercase tracking-wide hover:brightness-110\">Add</button> <button class=\"px-3 py-2 text-[var(--color-muted-text)] text-xs hover:text-[var(--color-white)] font-[var(--font-ui)]\">Cancel</button></div>"), _a = /* @__PURE__ */ J("<button class=\"mt-6 w-full py-3 border border-dashed border-white/[0.08] text-sm text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] hover:border-[var(--color-sunburst)]/30 flex items-center justify-center gap-2 font-[var(--font-ui)] uppercase tracking-wide transition-colors\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Section</button>"), va = /* @__PURE__ */ J("<div class=\"estimate-builder pb-20\"><div class=\"sticky top-0 z-10 border-b border-white/[0.06]\" style=\"background: var(--color-granite);\"><div class=\"flex items-center justify-between px-5 py-2\"><div class=\"flex items-center gap-4\"><button class=\"text-xs px-2 py-1 bg-white/[0.06] hover:bg-white/[0.1] disabled:opacity-30 disabled:cursor-not-allowed transition-colors flex items-center gap-1 font-[var(--font-ui)] text-[var(--color-concrete)]\" title=\"Undo (Ctrl+Z)\"><svg class=\"w-3.5 h-3.5\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M3 10h10a5 5 0 015 5v2M3 10l4-4M3 10l4 4\"></path></svg> Undo</button> <!></div> <div class=\"flex items-center gap-3 text-sm\"><span class=\"text-[var(--color-sage)] font-[var(--font-ui)] uppercase tracking-wider text-xs font-semibold\">Markup</span> <!></div></div></div> <div class=\"px-5 pt-6\"><!> <!></div> <!></div>");
function ya(e, t) {
	Fe(t, !0);
	let n = /* @__PURE__ */ I(null), r = /* @__PURE__ */ I(null), i = /* @__PURE__ */ I(!0), a = /* @__PURE__ */ I(!1), o = /* @__PURE__ */ I(""), s = la(t.projectId), c = da();
	async function l() {
		try {
			let e = await fetch(`/api/estimate/${t.projectId}`);
			if (!e.ok) throw Error(`Failed to load estimate: ${e.status}`);
			L(n, await e.json(), !0), s.register(() => K(n));
		} catch (e) {
			L(r, e.message, !0);
		} finally {
			L(i, !1);
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
			(e.ctrlKey || e.metaKey) && e.key === "z" && !e.shiftKey && (e.preventDefault(), K(n) && c.undo(K(n)) && u()), (e.ctrlKey || e.metaKey) && e.key === "s" && (e.preventDefault(), s.status === "dirty" && s.save());
		}
		return window.addEventListener("beforeunload", e), document.addEventListener("click", t, !0), window.addEventListener("keydown", r), () => {
			window.removeEventListener("beforeunload", e), document.removeEventListener("click", t, !0), window.removeEventListener("keydown", r), s.destroy();
		};
	});
	function u() {
		s.markDirty();
	}
	function d() {
		K(n) && c.snapshot(K(n));
	}
	function f() {
		d();
		let e = K(o).trim() || "New Section";
		K(n).sections.push({
			id: Si(),
			name: e,
			sort_order: K(n).sections.length,
			subcategories: []
		}), L(o, ""), L(a, !1), u();
	}
	function p(e) {
		d();
		let t = K(n).sections.findIndex((t) => t.id === e);
		t !== -1 && (K(n).sections.splice(t, 1), u());
	}
	function m(e, t) {
		let r = parseFloat(t.target.value);
		!isNaN(r) && r >= 0 && (K(n).globals[`${e}_markup`] = r, u());
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
	var g = yr(), _ = Qt(g), v = (e) => {
		Y(e, fa());
	}, y = (e) => {
		var t = pa(), n = z(t), i = z(n, !0);
		A(n), A(t), V(() => X(i, K(r))), Y(e, t);
	}, b = (e) => {
		var t = va(), r = z(t), i = z(r), l = z(i), g = z(l);
		ca(B(g, 2), {
			get status() {
				return s.status;
			},
			get savedAt() {
				return s.savedAt;
			},
			onsave: () => s.save()
		}), A(l);
		var _ = B(l, 2);
		kr(B(z(_), 2), 17, () => h, Tr, (e, t) => {
			var r = ma(), i = z(r), a = z(i, !0);
			A(i);
			var o = B(i, 2);
			Q(o), Ee(2), A(r), V(() => {
				X(a, K(t).label), Gr(o, K(n).globals[`${K(t).key}_markup`]);
			}), q("input", o, (e) => m(K(t).key, e)), Y(e, r);
		}), A(_), A(i), A(r);
		var v = B(r, 2), y = z(v), b = (e) => {
			var t = ha(), n = B(z(t), 4);
			A(t), q("click", n, () => L(a, !0)), Y(e, t);
		}, x = (e) => {
			var t = yr();
			kr(Qt(t), 17, () => K(n).sections, (e) => e.id, (e, t) => {
				ea(e, {
					get section() {
						return K(t);
					},
					get globals() {
						return K(n).globals;
					},
					onchange: u,
					onsnapshot: d,
					ondelete: p,
					get materialsDb() {
						return K(n).materials_db;
					},
					get ratesDb() {
						return K(n).rates_db;
					},
					get subcontractorsDb() {
						return K(n).subcontractors_db;
					}
				});
			}), Y(e, t);
		};
		Z(y, (e) => {
			K(n).sections.length === 0 && !K(a) ? e(b) : e(x, -1);
		});
		var S = B(y, 2), ee = (e) => {
			var t = ga(), n = z(t);
			Q(n);
			var r = B(n, 2), i = B(r, 2);
			A(t), q("keydown", n, (e) => {
				e.key === "Enter" && f(), e.key === "Escape" && L(a, !1);
			}), Zr(n, () => K(o), (e) => L(o, e)), q("click", r, f), q("click", i, () => L(a, !1)), Y(e, t);
		}, te = (e) => {
			var t = _a();
			q("click", t, () => L(a, !0)), Y(e, t);
		};
		Z(S, (e) => {
			K(a) ? e(ee) : K(n).sections.length > 0 && e(te, 1);
		}), A(v), ia(B(v, 2), { get estimate() {
			return K(n);
		} }), A(t), V(() => g.disabled = !c.canUndo), q("click", g, () => {
			c.undo(K(n)) && u();
		}), Y(e, t);
	};
	Z(_, (e) => {
		K(i) ? e(v) : K(r) ? e(y, 1) : K(n) && e(b, 2);
	}), Y(e, g), Ie();
}
fr([
	"click",
	"input",
	"keydown"
]);
//#endregion
//#region src/main.js
var ba = document.getElementById("estimate-root");
if (ba) {
	let e = ba.dataset.projectId;
	br(ya, {
		target: ba,
		props: { projectId: e }
	});
}
//#endregion
