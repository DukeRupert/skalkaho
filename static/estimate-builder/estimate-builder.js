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
var m = 1024, h = 2048, g = 4096, _ = 8192, v = 16384, y = 32768, b = 1 << 25, x = 65536, S = 1 << 19, ee = 1 << 20, C = 1 << 25, w = 65536, te = 1 << 21, T = 1 << 22, E = 1 << 23, D = Symbol("$state"), ne = Symbol("legacy props"), re = Symbol(""), O = new class extends Error {
	name = "StaleReactionError";
	message = "The reaction that called `getAbortSignal()` was re-run or destroyed";
}(), ie = !!globalThis.document?.contentType && /* @__PURE__ */ globalThis.document.contentType.includes("xml");
//#endregion
//#region node_modules/svelte/src/internal/client/errors.js
function ae() {
	throw Error("https://svelte.dev/e/async_derived_orphan");
}
function oe(e, t, n) {
	throw Error("https://svelte.dev/e/each_key_duplicate");
}
function se(e) {
	throw Error("https://svelte.dev/e/effect_in_teardown");
}
function ce() {
	throw Error("https://svelte.dev/e/effect_in_unowned_derived");
}
function le(e) {
	throw Error("https://svelte.dev/e/effect_orphan");
}
function ue() {
	throw Error("https://svelte.dev/e/effect_update_depth_exceeded");
}
function de(e) {
	throw Error("https://svelte.dev/e/props_invalid_value");
}
function fe() {
	throw Error("https://svelte.dev/e/state_descriptors_fixed");
}
function pe() {
	throw Error("https://svelte.dev/e/state_prototype_fixed");
}
function me() {
	throw Error("https://svelte.dev/e/state_unsafe_mutation");
}
function he() {
	throw Error("https://svelte.dev/e/svelte_boundary_reset_onerror");
}
//#endregion
//#region node_modules/svelte/src/constants.js
var ge = {}, k = Symbol(), _e = "http://www.w3.org/1999/xhtml";
function ve(e) {
	console.warn("https://svelte.dev/e/hydration_mismatch");
}
function ye() {
	console.warn("https://svelte.dev/e/select_multiple_invalid_value");
}
function be() {
	console.warn("https://svelte.dev/e/svelte_boundary_reset_noop");
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/hydration.js
var A = !1;
function xe(e) {
	A = e;
}
var j;
function Se(e) {
	if (e === null) throw ve(), ge;
	return j = e;
}
function Ce() {
	return Se(/* @__PURE__ */ Xt(j));
}
function M(e) {
	if (A) {
		if (/* @__PURE__ */ Xt(j) !== null) throw ve(), ge;
		j = e;
	}
}
function we(e = 1) {
	if (A) {
		for (var t = e, n = j; t--;) n = /* @__PURE__ */ Xt(n);
		j = n;
	}
}
function Te(e = !0) {
	for (var t = 0, n = j;;) {
		if (n.nodeType === 8) {
			var r = n.data;
			if (r === "]") {
				if (t === 0) return n;
				--t;
			} else (r === "[" || r === "[!" || r[0] === "[" && !isNaN(Number(r.slice(1)))) && (t += 1);
		}
		var i = /* @__PURE__ */ Xt(n);
		e && n.remove(), n = i;
	}
}
function Ee(e) {
	if (!e || e.nodeType !== 8) throw ve(), ge;
	return e.data;
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/equality.js
function De(e) {
	return e === this.v;
}
function Oe(e, t) {
	return e == e ? e !== t || typeof e == "object" && !!e || typeof e == "function" : t == t;
}
function ke(e) {
	return !Oe(e, this.v);
}
//#endregion
//#region node_modules/svelte/src/internal/flags/index.js
var Ae = !1, je = !1, N = null;
function Me(e) {
	N = e;
}
function Ne(e, t = !1, n) {
	N = {
		p: N,
		i: !1,
		c: null,
		e: null,
		s: e,
		x: null,
		r: G,
		l: je && !t ? {
			s: null,
			u: null,
			$: []
		} : null
	};
}
function Pe(e) {
	var t = N, n = t.e;
	if (n !== null) {
		t.e = null;
		for (var r of n) mn(r);
	}
	return e !== void 0 && (t.x = e), t.i = !0, N = t.p, e ?? {};
}
function Fe() {
	return !je || N !== null && N.l === null;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/task.js
var Ie = [];
function Le() {
	var e = Ie;
	Ie = [], f(e);
}
function Re(e) {
	if (Ie.length === 0 && !Qe) {
		var t = Ie;
		queueMicrotask(() => {
			t === Ie && Le();
		});
	}
	Ie.push(e);
}
function ze() {
	for (; Ie.length > 0;) Le();
}
function Be(e) {
	var t = G;
	if (t === null) return W.f |= E, e;
	if (!(t.f & 32768) && !(t.f & 4)) throw e;
	Ve(e, t);
}
function Ve(e, t) {
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
var He = ~(h | g | m);
function P(e, t) {
	e.f = e.f & He | t;
}
function Ue(e) {
	e.f & 512 || e.deps === null ? P(e, m) : P(e, g);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/utils.js
function We(e) {
	if (e !== null) for (let t of e) !(t.f & 2) || !(t.f & 65536) || (t.f ^= w, We(t.deps));
}
function Ge(e, t, n) {
	e.f & 2048 ? t.add(e) : e.f & 4096 && n.add(e), We(e.deps), P(e, m);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/store.js
var Ke = !1, qe = !1;
function Je(e) {
	var t = qe;
	try {
		return qe = !1, [e(), qe];
	} finally {
		qe = t;
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/batch.js
var Ye = /* @__PURE__ */ new Set(), F = null, Xe = null, I = null, Ze = null, Qe = !1, $e = !1, et = null, tt = null, nt = 0, rt = 1, it = class e {
	id = rt++;
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
			for (var n of t.d) P(n, h), this.schedule(n);
			for (n of t.m) P(n, g), this.schedule(n);
		}
	}
	#d() {
		nt++ > 1e3 && ot();
		let t = this.#a;
		this.#a = [], this.apply();
		var n = et = [], r = [], i = tt = [];
		for (let e of t) try {
			this.#f(e, n, r);
		} catch (t) {
			throw pt(e), t;
		}
		if (F = null, i.length > 0) {
			var a = e.ensure();
			for (let e of i) a.schedule(e);
		}
		if (et = null, tt = null, this.#u()) {
			this.#p(r), this.#p(n);
			for (let [e, t] of this.#c) ft(e, t);
		} else {
			this.#n === 0 && Ye.delete(this), this.#o.clear(), this.#s.clear();
			for (let e of this.#e) e(this);
			this.#e.clear(), Xe = this, ct(r), ct(n), Xe = null, this.#i?.resolve();
		}
		var o = F;
		if (this.#a.length > 0) {
			let e = o ??= this;
			e.#a.push(...this.#a.filter((t) => !e.#a.includes(t)));
		}
		o !== null && (Ye.add(o), o.#d()), Ye.has(this) || this.#m();
	}
	#f(e, t, n) {
		e.f ^= m;
		for (var r = e.first; r !== null;) {
			var i = r.f, a = (i & 96) != 0;
			if (!(a && i & 1024 || i & 8192 || this.#c.has(r)) && r.fn !== null) {
				a ? r.f ^= m : i & 4 ? t.push(r) : Ae && i & 16777224 ? n.push(r) : Jn(r) && (i & 16 && this.#s.add(r), $n(r));
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
		for (var t = 0; t < e.length; t += 1) Ge(e[t], this.#o, this.#s);
	}
	capture(e, t) {
		t !== k && !this.previous.has(e) && this.previous.set(e, t), e.f & 8388608 || (this.current.set(e, e.v), I?.set(e, e.v));
	}
	activate() {
		F = this;
	}
	deactivate() {
		F = null, I = null;
	}
	flush() {
		try {
			if ($e = !0, F = this, !this.#u()) {
				for (let e of this.#o) this.#s.delete(e), P(e, h), this.schedule(e);
				for (let e of this.#s) P(e, g), this.schedule(e);
			}
			this.#d();
		} finally {
			nt = 0, Ze = null, et = null, tt = null, $e = !1, F = null, I = null, Mt.clear();
		}
	}
	discard() {
		for (let e of this.#t) e(this);
		this.#t.clear();
	}
	#m() {
		for (let s of Ye) {
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
					for (var a of t) lt(a, n, r, i);
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
		--this.#n, e && --this.#r, !(this.#l || t) && (this.#l = !0, Re(() => {
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
			$e || (Ye.add(F), Qe || Re(() => {
				F === t && t.flush();
			}));
		}
		return F;
	}
	apply() {
		if (!Ae || !this.is_fork && Ye.size === 1) {
			I = null;
			return;
		}
		I = new Map(this.current);
		for (let e of Ye) if (e !== this) for (let [t, n] of e.previous) I.has(t) || I.set(t, n);
	}
	schedule(e) {
		if (Ze = e, e.b?.is_pending && e.f & 16777228 && !(e.f & 32768)) {
			e.b.defer_effect(e);
			return;
		}
		for (var t = e; t.parent !== null;) {
			t = t.parent;
			var n = t.f;
			if (et !== null && t === G && (Ae || (W === null || !(W.f & 2)) && !Ke)) return;
			if (n & 96) {
				if (!(n & 1024)) return;
				t.f ^= m;
			}
		}
		this.#a.push(t);
	}
};
function at(e) {
	var t = Qe;
	Qe = !0;
	try {
		var n;
		for (e && (F !== null && !F.is_fork && F.flush(), n = e());;) {
			if (ze(), F === null) return n;
			F.flush();
		}
	} finally {
		Qe = t;
	}
}
function ot() {
	try {
		ue();
	} catch (e) {
		Ve(e, Ze);
	}
}
var st = null;
function ct(e) {
	var t = e.length;
	if (t !== 0) {
		for (var n = 0; n < t;) {
			var r = e[n++];
			if (!(r.f & 24576) && Jn(r) && (st = /* @__PURE__ */ new Set(), $n(r), r.deps === null && r.first === null && r.nodes === null && r.teardown === null && r.ac === null && Tn(r), st?.size > 0)) {
				Mt.clear();
				for (let e of st) {
					if (e.f & 24576) continue;
					let t = [e], n = e.parent;
					for (; n !== null;) st.has(n) && (st.delete(n), t.push(n)), n = n.parent;
					for (let e = t.length - 1; e >= 0; e--) {
						let n = t[e];
						n.f & 24576 || $n(n);
					}
				}
				st.clear();
			}
		}
		st = null;
	}
}
function lt(e, t, n, r) {
	if (!n.has(e) && (n.add(e), e.reactions !== null)) for (let i of e.reactions) {
		let e = i.f;
		e & 2 ? lt(i, t, n, r) : e & 4194320 && !(e & 2048) && ut(i, t, r) && (P(i, h), dt(i));
	}
}
function ut(e, t, r) {
	let i = r.get(e);
	if (i !== void 0) return i;
	if (e.deps !== null) for (let i of e.deps) {
		if (n.call(t, i)) return !0;
		if (i.f & 2 && ut(i, t, r)) return r.set(i, !0), !0;
	}
	return r.set(e, !1), !1;
}
function dt(e) {
	F.schedule(e);
}
function ft(e, t) {
	if (!(e.f & 32 && e.f & 1024)) {
		e.f & 2048 ? t.d.push(e) : e.f & 4096 && t.m.push(e), P(e, m);
		for (var n = e.first; n !== null;) ft(n, t), n = n.next;
	}
}
function pt(e) {
	P(e, m);
	for (var t = e.first; t !== null;) pt(t), t = t.next;
}
//#endregion
//#region node_modules/svelte/src/reactivity/create-subscriber.js
function mt(e) {
	let t = 0, n = Pt(0), r;
	return () => {
		dn() && (q(n), vn(() => (t === 0 && (r = rr(() => e(() => Rt(n)))), t += 1, () => {
			Re(() => {
				--t, t === 0 && (r?.(), r = void 0, Rt(n));
			});
		})));
	};
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/boundary.js
var ht = x | S;
function gt(e, t, n, r) {
	new _t(e, t, n, r);
}
var _t = class {
	parent;
	is_pending = !1;
	transform_error;
	#e;
	#t = A ? j : null;
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
	#h = mt(() => (this.#m = Pt(this.#l), () => {
		this.#m = null;
	}));
	constructor(e, t, n, r) {
		this.#e = e, this.#n = t, this.#r = (e) => {
			var t = G;
			t.b = this, t.f |= 128, n(e);
		}, this.parent = G.b, this.transform_error = r ?? this.parent?.transform_error ?? ((e) => e), this.#i = yn(() => {
			if (A) {
				let e = this.#t;
				Ce();
				let t = e.data === "[!";
				if (e.data.startsWith("[?")) {
					let t = JSON.parse(e.data.slice(2));
					this.#_(t);
				} else t ? this.#v() : this.#g();
			} else this.#y();
		}, ht), A && (this.#e = j);
	}
	#g() {
		try {
			this.#a = bn(() => this.#r(this.#e));
		} catch (e) {
			this.error(e);
		}
	}
	#_(e) {
		let t = this.#n.failed;
		t && (this.#s = bn(() => {
			t(this.#e, () => e, () => () => {});
		}));
	}
	#v() {
		let e = this.#n.pending;
		e && (this.is_pending = !0, this.#o = bn(() => e(this.#e)), Re(() => {
			var e = this.#c = document.createDocumentFragment(), t = Jt();
			e.append(t), this.#a = this.#x(() => bn(() => this.#r(t))), this.#u === 0 && (this.#e.before(e), this.#c = null, En(this.#o, () => {
				this.#o = null;
			}), this.#b(F));
		}));
	}
	#y() {
		try {
			if (this.is_pending = this.has_pending_snippet(), this.#u = 0, this.#l = 0, this.#a = bn(() => {
				this.#r(this.#e);
			}), this.#u > 0) {
				var e = this.#c = document.createDocumentFragment();
				An(this.#a, e);
				let t = this.#n.pending;
				this.#o = bn(() => t(this.#e));
			} else this.#b(F);
		} catch (e) {
			this.error(e);
		}
	}
	#b(e) {
		this.is_pending = !1;
		for (let t of this.#f) P(t, h), e.schedule(t);
		for (let t of this.#p) P(t, g), e.schedule(t);
		this.#f.clear(), this.#p.clear();
	}
	defer_effect(e) {
		Ge(e, this.#f, this.#p);
	}
	is_rendered() {
		return !this.is_pending && (!this.parent || this.parent.is_rendered());
	}
	has_pending_snippet() {
		return !!this.#n.pending;
	}
	#x(e) {
		var t = G, n = W, r = N;
		Ln(this.#i), In(this.#i), Me(this.#i.ctx);
		try {
			return it.ensure(), e();
		} catch (e) {
			return Be(e), null;
		} finally {
			Ln(t), In(n), Me(r);
		}
	}
	#S(e, t) {
		if (!this.has_pending_snippet()) {
			this.parent && this.parent.#S(e, t);
			return;
		}
		this.#u += e, this.#u === 0 && (this.#b(t), this.#o && En(this.#o, () => {
			this.#o = null;
		}), this.#c &&= (this.#e.before(this.#c), null));
	}
	update_pending_count(e, t) {
		this.#S(e, t), this.#l += e, !(!this.#m || this.#d) && (this.#d = !0, Re(() => {
			this.#d = !1, this.#m && It(this.#m, this.#l);
		}));
	}
	get_effect_pending() {
		return this.#h(), q(this.#m);
	}
	error(e) {
		var t = this.#n.onerror;
		let n = this.#n.failed;
		if (!t && !n) throw e;
		this.#a &&= (U(this.#a), null), this.#o &&= (U(this.#o), null), this.#s &&= (U(this.#s), null), A && (Se(this.#t), we(), Se(Te()));
		var r = !1, i = !1;
		let a = () => {
			if (r) {
				be();
				return;
			}
			r = !0, i && he(), this.#s !== null && En(this.#s, () => {
				this.#s = null;
			}), this.#x(() => {
				this.#y();
			});
		}, o = (e) => {
			try {
				i = !0, t?.(e, a), i = !1;
			} catch (e) {
				Ve(e, this.#i && this.#i.parent);
			}
			n && (this.#s = this.#x(() => {
				try {
					return bn(() => {
						var t = G;
						t.b = this, t.f |= 128, n(this.#e, () => e, () => a);
					});
				} catch (e) {
					return Ve(e, this.#i.parent), null;
				}
			}));
		};
		Re(() => {
			var t;
			try {
				t = this.transform_error(e);
			} catch (e) {
				Ve(e, this.#i && this.#i.parent);
				return;
			}
			typeof t == "object" && t && typeof t.then == "function" ? t.then(o, (e) => Ve(e, this.#i && this.#i.parent)) : o(t);
		});
	}
};
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/async.js
function vt(e, t, n, r) {
	let i = Fe() ? St : wt;
	var a = e.filter((e) => !e.settled);
	if (n.length === 0 && a.length === 0) {
		r(t.map(i));
		return;
	}
	var o = G, s = yt(), c = a.length === 1 ? a[0].promise : a.length > 1 ? Promise.all(a.map((e) => e.promise)) : null;
	function l(e) {
		s();
		try {
			r(e);
		} catch (e) {
			o.f & 16384 || Ve(e, o);
		}
		bt();
	}
	if (n.length === 0) {
		c.then(() => l(t.map(i)));
		return;
	}
	var u = xt();
	function d() {
		Promise.all(n.map((e) => /* @__PURE__ */ Ct(e))).then((e) => l([...t.map(i), ...e])).catch((e) => Ve(e, o)).finally(() => u());
	}
	c ? c.then(() => {
		s(), d(), bt();
	}) : d();
}
function yt() {
	var e = G, t = W, n = N, r = F;
	return function(i = !0) {
		Ln(e), In(t), Me(n), i && !(e.f & 16384) && (r?.activate(), r?.apply());
	};
}
function bt(e = !0) {
	Ln(null), In(null), Me(null), e && F?.deactivate();
}
function xt() {
	var e = G.b, t = F, n = e.is_rendered();
	return e.update_pending_count(1, t), t.increment(n), (r = !1) => {
		e.update_pending_count(-1, t), t.decrement(n, r);
	};
}
/* @__NO_SIDE_EFFECTS__ */
function St(e) {
	var t = 2 | h, n = W !== null && W.f & 2 ? W : null;
	return G !== null && (G.f |= S), {
		ctx: N,
		deps: null,
		effects: null,
		equals: De,
		f: t,
		fn: e,
		reactions: null,
		rv: 0,
		v: k,
		wv: 0,
		parent: n ?? G,
		ac: null
	};
}
/* @__NO_SIDE_EFFECTS__ */
function Ct(e, t, n) {
	let r = G;
	r === null && ae();
	var i = void 0, a = Pt(k), o = !W, s = /* @__PURE__ */ new Map();
	return _n(() => {
		var t = G, n = p();
		i = n.promise;
		try {
			Promise.resolve(e()).then(n.resolve, n.reject).finally(bt);
		} catch (e) {
			n.reject(e), bt();
		}
		var c = F;
		if (o) {
			if (t.f & 32768) var l = xt();
			if (r.b.is_rendered()) s.get(c)?.reject(O), s.delete(c);
			else {
				for (let e of s.values()) e.reject(O);
				s.clear();
			}
			s.set(c, n);
		}
		let u = (e, n = void 0) => {
			if (l && l(n === O), !(n === O || t.f & 16384)) {
				if (c.activate(), n) a.f |= E, It(a, n);
				else {
					a.f & 8388608 && (a.f ^= E), It(a, e);
					for (let [e, t] of s) {
						if (s.delete(e), e === c) break;
						t.reject(O);
					}
				}
				c.deactivate();
			}
		};
		n.promise.then(u, (e) => u(null, e || "unknown"));
	}), fn(() => {
		for (let e of s.values()) e.reject(O);
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
function L(e) {
	let t = /* @__PURE__ */ St(e);
	return Ae || zn(t), t;
}
/* @__NO_SIDE_EFFECTS__ */
function wt(e) {
	let t = /* @__PURE__ */ St(e);
	return t.equals = ke, t;
}
function Tt(e) {
	var t = e.effects;
	if (t !== null) {
		e.effects = null;
		for (var n = 0; n < t.length; n += 1) U(t[n]);
	}
}
function Et(e) {
	for (var t = e.parent; t !== null;) {
		if (!(t.f & 2)) return t.f & 16384 ? null : t;
		t = t.parent;
	}
	return null;
}
function Dt(e) {
	var t, n = G;
	Ln(Et(e));
	try {
		e.f &= ~w, Tt(e), t = Xn(e);
	} finally {
		Ln(n);
	}
	return t;
}
function Ot(e) {
	var t = Dt(e);
	if (!e.equals(t) && (e.wv = qn(), (!F?.is_fork || e.deps === null) && (e.v = t, e.deps === null))) {
		P(e, m);
		return;
	}
	Nn || (I === null ? Ue(e) : (dn() || F?.is_fork) && I.set(e, t));
}
function kt(e) {
	if (e.effects !== null) for (let t of e.effects) (t.teardown || t.ac) && (t.teardown?.(), t.ac?.abort(O), t.teardown = d, t.ac = null, Qn(t, 0), Sn(t));
}
function At(e) {
	if (e.effects !== null) for (let t of e.effects) t.teardown && $n(t);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/sources.js
var jt = /* @__PURE__ */ new Set(), Mt = /* @__PURE__ */ new Map(), Nt = !1;
function Pt(e, t) {
	return {
		f: 0,
		v: e,
		reactions: null,
		equals: De,
		rv: 0,
		wv: 0
	};
}
/* @__NO_SIDE_EFFECTS__ */
function R(e, t) {
	let n = Pt(e, t);
	return zn(n), n;
}
/* @__NO_SIDE_EFFECTS__ */
function Ft(e, t = !1, n = !0) {
	let r = Pt(e);
	return t || (r.equals = ke), je && n && N !== null && N.l !== null && (N.l.s ??= []).push(r), r;
}
function z(e, t, r = !1) {
	return W !== null && (!Fn || W.f & 131072) && Fe() && W.f & 4325394 && (Rn === null || !n.call(Rn, e)) && me(), It(e, r ? Bt(t) : t, tt);
}
function It(e, t, n = null) {
	if (!e.equals(t)) {
		var r = e.v;
		Nn ? Mt.set(e, t) : Mt.set(e, r), e.v = t;
		var i = it.ensure();
		if (i.capture(e, r), e.f & 2) {
			let t = e;
			e.f & 2048 && Dt(t), Ue(t);
		}
		e.wv = qn(), zt(e, h, n), Fe() && G !== null && G.f & 1024 && !(G.f & 96) && (Vn === null ? Hn([e]) : Vn.push(e)), !i.is_fork && jt.size > 0 && !Nt && Lt();
	}
	return t;
}
function Lt() {
	Nt = !1;
	for (let e of jt) e.f & 1024 && P(e, g), Jn(e) && $n(e);
	jt.clear();
}
function Rt(e) {
	z(e, e.v + 1);
}
function zt(e, t, n) {
	var r = e.reactions;
	if (r !== null) for (var i = Fe(), a = r.length, o = 0; o < a; o++) {
		var s = r[o], c = s.f;
		if (!(!i && s === G)) {
			var l = (c & h) === 0;
			if (l && P(s, t), c & 2) {
				var u = s;
				I?.delete(u), c & 65536 || (c & 512 && (s.f |= w), zt(u, g, n));
			} else if (l) {
				var d = s;
				c & 16 && st !== null && st.add(d), n === null ? dt(d) : n.push(d);
			}
		}
	}
}
function Bt(t) {
	if (typeof t != "object" || !t || D in t) return t;
	let n = l(t);
	if (n !== s && n !== c) return t;
	var r = /* @__PURE__ */ new Map(), i = e(t), o = /* @__PURE__ */ R(0), u = null, d = Gn, f = (e) => {
		if (Gn === d) return e();
		var t = W, n = Gn;
		In(null), Kn(d);
		var r = e();
		return In(t), Kn(n), r;
	};
	return i && r.set("length", /* @__PURE__ */ R(t.length, u)), new Proxy(t, {
		defineProperty(e, t, n) {
			(!("value" in n) || n.configurable === !1 || n.enumerable === !1 || n.writable === !1) && fe();
			var i = r.get(t);
			return i === void 0 ? f(() => {
				var e = /* @__PURE__ */ R(n.value, u);
				return r.set(t, e), e;
			}) : z(i, n.value, !0), !0;
		},
		deleteProperty(e, t) {
			var n = r.get(t);
			if (n === void 0) {
				if (t in e) {
					let e = f(() => /* @__PURE__ */ R(k, u));
					r.set(t, e), Rt(o);
				}
			} else z(n, k), Rt(o);
			return !0;
		},
		get(e, n, i) {
			if (n === D) return t;
			var o = r.get(n), s = n in e;
			if (o === void 0 && (!s || a(e, n)?.writable) && (o = f(() => /* @__PURE__ */ R(Bt(s ? e[n] : k), u)), r.set(n, o)), o !== void 0) {
				var c = q(o);
				return c === k ? void 0 : c;
			}
			return Reflect.get(e, n, i);
		},
		getOwnPropertyDescriptor(e, t) {
			var n = Reflect.getOwnPropertyDescriptor(e, t);
			if (n && "value" in n) {
				var i = r.get(t);
				i && (n.value = q(i));
			} else if (n === void 0) {
				var a = r.get(t), o = a?.v;
				if (a !== void 0 && o !== k) return {
					enumerable: !0,
					configurable: !0,
					value: o,
					writable: !0
				};
			}
			return n;
		},
		has(e, t) {
			if (t === D) return !0;
			var n = r.get(t), i = n !== void 0 && n.v !== k || Reflect.has(e, t);
			return (n !== void 0 || G !== null && (!i || a(e, t)?.writable)) && (n === void 0 && (n = f(() => /* @__PURE__ */ R(i ? Bt(e[t]) : k, u)), r.set(t, n)), q(n) === k) ? !1 : i;
		},
		set(e, t, n, s) {
			var c = r.get(t), l = t in e;
			if (i && t === "length") for (var d = n; d < c.v; d += 1) {
				var p = r.get(d + "");
				p === void 0 ? d in e && (p = f(() => /* @__PURE__ */ R(k, u)), r.set(d + "", p)) : z(p, k);
			}
			if (c === void 0) (!l || a(e, t)?.writable) && (c = f(() => /* @__PURE__ */ R(void 0, u)), z(c, Bt(n)), r.set(t, c));
			else {
				l = c.v !== k;
				var m = f(() => Bt(n));
				z(c, m);
			}
			var h = Reflect.getOwnPropertyDescriptor(e, t);
			if (h?.set && h.set.call(s, n), !l) {
				if (i && typeof t == "string") {
					var g = r.get("length"), _ = Number(t);
					Number.isInteger(_) && _ >= g.v && z(g, _ + 1);
				}
				Rt(o);
			}
			return !0;
		},
		ownKeys(e) {
			q(o);
			var t = Reflect.ownKeys(e).filter((e) => {
				var t = r.get(e);
				return t === void 0 || t.v !== k;
			});
			for (var [n, i] of r) i.v !== k && !(n in e) && t.push(n);
			return t;
		},
		setPrototypeOf() {
			pe();
		}
	});
}
function Vt(e) {
	try {
		if (typeof e == "object" && e && D in e) return e[D];
	} catch {}
	return e;
}
function Ht(e, t) {
	return Object.is(Vt(e), Vt(t));
}
var Ut, Wt, Gt, Kt;
function qt() {
	if (Ut === void 0) {
		Ut = window, document, Wt = /Firefox/.test(navigator.userAgent);
		var e = Element.prototype, t = Node.prototype, n = Text.prototype;
		Gt = a(t, "firstChild").get, Kt = a(t, "nextSibling").get, u(e) && (e.__click = void 0, e.__className = void 0, e.__attributes = null, e.__style = void 0, e.__e = void 0), u(n) && (n.__t = void 0);
	}
}
function Jt(e = "") {
	return document.createTextNode(e);
}
/* @__NO_SIDE_EFFECTS__ */
function Yt(e) {
	return Gt.call(e);
}
/* @__NO_SIDE_EFFECTS__ */
function Xt(e) {
	return Kt.call(e);
}
function B(e, t) {
	if (!A) return /* @__PURE__ */ Yt(e);
	var n = /* @__PURE__ */ Yt(j);
	if (n === null) n = j.appendChild(Jt());
	else if (t && n.nodeType !== 3) {
		var r = Jt();
		return n?.before(r), Se(r), r;
	}
	return t && tn(n), Se(n), n;
}
function Zt(e, t = !1) {
	if (!A) {
		var n = /* @__PURE__ */ Yt(e);
		return n instanceof Comment && n.data === "" ? /* @__PURE__ */ Xt(n) : n;
	}
	if (t) {
		if (j?.nodeType !== 3) {
			var r = Jt();
			return j?.before(r), Se(r), r;
		}
		tn(j);
	}
	return j;
}
function V(e, t = 1, n = !1) {
	let r = A ? j : e;
	for (var i; t--;) i = r, r = /* @__PURE__ */ Xt(r);
	if (!A) return r;
	if (n) {
		if (r?.nodeType !== 3) {
			var a = Jt();
			return r === null ? i?.after(a) : r.before(a), Se(a), a;
		}
		tn(r);
	}
	return Se(r), r;
}
function Qt(e) {
	e.textContent = "";
}
function $t() {
	return !Ae || st !== null ? !1 : (G.f & y) !== 0;
}
function en(e, t, n) {
	let r = n ? { is: n } : void 0;
	return document.createElementNS(t ?? "http://www.w3.org/1999/xhtml", e, r);
}
function tn(e) {
	if (e.nodeValue.length < 65536) return;
	let t = e.nextSibling;
	for (; t !== null && t.nodeType === 3;) t.remove(), e.nodeValue += t.nodeValue, t = e.nextSibling;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/misc.js
function nn(e, t) {
	if (t) {
		let t = document.body;
		e.autofocus = !0, Re(() => {
			document.activeElement === t && e.focus();
		});
	}
}
var rn = !1;
function an() {
	rn || (rn = !0, document.addEventListener("reset", (e) => {
		Promise.resolve().then(() => {
			if (!e.defaultPrevented) for (let t of e.target.elements) t.__on_r?.();
		});
	}, { capture: !0 }));
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/shared.js
function on(e) {
	var t = W, n = G;
	In(null), Ln(null);
	try {
		return e();
	} finally {
		In(t), Ln(n);
	}
}
function sn(e, t, n, r = n) {
	e.addEventListener(t, () => on(n));
	let i = e.__on_r;
	i ? e.__on_r = () => {
		i(), r(!0);
	} : e.__on_r = () => r(!0), an();
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/effects.js
function cn(e) {
	G === null && (W === null && le(e), ce()), Nn && se(e);
}
function ln(e, t) {
	var n = t.last;
	n === null ? t.last = t.first = e : (n.next = e, e.prev = n, t.last = e);
}
function un(e, t) {
	var n = G;
	n !== null && n.f & 8192 && (e |= _);
	var r = {
		ctx: N,
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
	if (e & 4) et === null ? it.ensure().schedule(r) : et.push(r);
	else if (t !== null) {
		try {
			$n(r);
		} catch (e) {
			throw U(r), e;
		}
		i.deps === null && i.teardown === null && i.nodes === null && i.first === i.last && !(i.f & 524288) && (i = i.first, e & 16 && e & 65536 && i !== null && (i.f |= x));
	}
	if (i !== null && (i.parent = n, n !== null && ln(i, n), W !== null && W.f & 2 && !(e & 64))) {
		var a = W;
		(a.effects ??= []).push(i);
	}
	return r;
}
function dn() {
	return W !== null && !Fn;
}
function fn(e) {
	let t = un(8, null);
	return P(t, m), t.teardown = e, t;
}
function pn(e) {
	cn("$effect");
	var t = G.f;
	if (!W && t & 32 && !(t & 32768)) {
		var n = N;
		(n.e ??= []).push(e);
	} else return mn(e);
}
function mn(e) {
	return un(4 | ee, e);
}
function hn(e) {
	it.ensure();
	let t = un(64 | S, e);
	return (e = {}) => new Promise((n) => {
		e.outro ? En(t, () => {
			U(t), n(void 0);
		}) : (U(t), n(void 0));
	});
}
function gn(e) {
	return un(4, e);
}
function _n(e) {
	return un(T | S, e);
}
function vn(e, t = 0) {
	return un(8 | t, e);
}
function H(e, t = [], n = [], r = []) {
	vt(r, t, n, (t) => {
		un(8, () => e(...t.map(q)));
	});
}
function yn(e, t = 0) {
	return un(16 | t, e);
}
function bn(e) {
	return un(32 | S, e);
}
function xn(e) {
	var t = e.teardown;
	if (t !== null) {
		let e = Nn, n = W;
		Pn(!0), In(null);
		try {
			t.call(null);
		} finally {
			Pn(e), In(n);
		}
	}
}
function Sn(e, t = !1) {
	var n = e.first;
	for (e.first = e.last = null; n !== null;) {
		let e = n.ac;
		e !== null && on(() => {
			e.abort(O);
		});
		var r = n.next;
		n.f & 64 ? n.parent = null : U(n, t), n = r;
	}
}
function Cn(e) {
	for (var t = e.first; t !== null;) {
		var n = t.next;
		t.f & 32 || U(t), t = n;
	}
}
function U(e, t = !0) {
	var n = !1;
	(t || e.f & 262144) && e.nodes !== null && e.nodes.end !== null && (wn(e.nodes.start, e.nodes.end), n = !0), P(e, b), Sn(e, t && !n), Qn(e, 0);
	var r = e.nodes && e.nodes.t;
	if (r !== null) for (let e of r) e.stop();
	xn(e), e.f ^= b, e.f |= v;
	var i = e.parent;
	i !== null && i.first !== null && Tn(e), e.next = e.prev = e.teardown = e.ctx = e.deps = e.fn = e.nodes = e.ac = null;
}
function wn(e, t) {
	for (; e !== null;) {
		var n = e === t ? null : /* @__PURE__ */ Xt(e);
		e.remove(), e = n;
	}
}
function Tn(e) {
	var t = e.parent, n = e.prev, r = e.next;
	n !== null && (n.next = r), r !== null && (r.prev = n), t !== null && (t.first === e && (t.first = r), t.last === e && (t.last = n));
}
function En(e, t, n = !0) {
	var r = [];
	Dn(e, r, !0);
	var i = () => {
		n && U(e), t && t();
	}, a = r.length;
	if (a > 0) {
		var o = () => --a || i();
		for (var s of r) s.out(o);
	} else i();
}
function Dn(e, t, n) {
	if (!(e.f & 8192)) {
		e.f ^= _;
		var r = e.nodes && e.nodes.t;
		if (r !== null) for (let e of r) (e.is_global || n) && t.push(e);
		for (var i = e.first; i !== null;) {
			var a = i.next, o = (i.f & 65536) != 0 || (i.f & 32) != 0 && (e.f & 16) != 0;
			Dn(i, t, o ? n : !1), i = a;
		}
	}
}
function On(e) {
	kn(e, !0);
}
function kn(e, t) {
	if (e.f & 8192) {
		e.f ^= _, e.f & 1024 || (P(e, h), it.ensure().schedule(e));
		for (var n = e.first; n !== null;) {
			var r = n.next, i = (n.f & 65536) != 0 || (n.f & 32) != 0;
			kn(n, i ? t : !1), n = r;
		}
		var a = e.nodes && e.nodes.t;
		if (a !== null) for (let e of a) (e.is_global || t) && e.in();
	}
}
function An(e, t) {
	if (e.nodes) for (var n = e.nodes.start, r = e.nodes.end; n !== null;) {
		var i = n === r ? null : /* @__PURE__ */ Xt(n);
		t.append(n), n = i;
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/legacy.js
var jn = null, Mn = !1, Nn = !1;
function Pn(e) {
	Nn = e;
}
var W = null, Fn = !1;
function In(e) {
	W = e;
}
var G = null;
function Ln(e) {
	G = e;
}
var Rn = null;
function zn(e) {
	W !== null && (!Ae || W.f & 2) && (Rn === null ? Rn = [e] : Rn.push(e));
}
var K = null, Bn = 0, Vn = null;
function Hn(e) {
	Vn = e;
}
var Un = 1, Wn = 0, Gn = Wn;
function Kn(e) {
	Gn = e;
}
function qn() {
	return ++Un;
}
function Jn(e) {
	var t = e.f;
	if (t & 2048) return !0;
	if (t & 2 && (e.f &= ~w), t & 4096) {
		for (var n = e.deps, r = n.length, i = 0; i < r; i++) {
			var a = n[i];
			if (Jn(a) && Ot(a), a.wv > e.wv) return !0;
		}
		t & 512 && I === null && P(e, m);
	}
	return !1;
}
function Yn(e, t, r = !0) {
	var i = e.reactions;
	if (i !== null && !(!Ae && Rn !== null && n.call(Rn, e))) for (var a = 0; a < i.length; a++) {
		var o = i[a];
		o.f & 2 ? Yn(o, t, !1) : t === o && (r ? P(o, h) : o.f & 1024 && P(o, g), dt(o));
	}
}
function Xn(e) {
	var t = K, n = Bn, r = Vn, i = W, a = Rn, o = N, s = Fn, c = Gn, l = e.f;
	K = null, Bn = 0, Vn = null, W = l & 96 ? null : e, Rn = null, Me(e.ctx), Fn = !1, Gn = ++Wn, e.ac !== null && (on(() => {
		e.ac.abort(O);
	}), e.ac = null);
	try {
		e.f |= te;
		var u = e.fn, d = u();
		e.f |= y;
		var f = e.deps, p = F?.is_fork;
		if (K !== null) {
			var m;
			if (p || Qn(e, Bn), f !== null && Bn > 0) for (f.length = Bn + K.length, m = 0; m < K.length; m++) f[Bn + m] = K[m];
			else e.deps = f = K;
			if (dn() && e.f & 512) for (m = Bn; m < f.length; m++) (f[m].reactions ??= []).push(e);
		} else !p && f !== null && Bn < f.length && (Qn(e, Bn), f.length = Bn);
		if (Fe() && Vn !== null && !Fn && f !== null && !(e.f & 6146)) for (m = 0; m < Vn.length; m++) Yn(Vn[m], e);
		if (i !== null && i !== e) {
			if (Wn++, i.deps !== null) for (let e = 0; e < n; e += 1) i.deps[e].rv = Wn;
			if (t !== null) for (let e of t) e.rv = Wn;
			Vn !== null && (r === null ? r = Vn : r.push(...Vn));
		}
		return e.f & 8388608 && (e.f ^= E), d;
	} catch (e) {
		return Be(e);
	} finally {
		e.f ^= te, K = t, Bn = n, Vn = r, W = i, Rn = a, Me(o), Fn = s, Gn = c;
	}
}
function Zn(e, r) {
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
		s.f & 512 && (s.f ^= 512, s.f &= ~w), Ue(s), kt(s), Qn(s, 0);
	}
}
function Qn(e, t) {
	var n = e.deps;
	if (n !== null) for (var r = t; r < n.length; r++) Zn(e, n[r]);
}
function $n(e) {
	var t = e.f;
	if (!(t & 16384)) {
		P(e, m);
		var n = G, r = Mn;
		G = e, Mn = !0;
		try {
			t & 16777232 ? Cn(e) : Sn(e), xn(e);
			var i = Xn(e);
			e.teardown = typeof i == "function" ? i : null, e.wv = Un;
		} finally {
			Mn = r, G = n;
		}
	}
}
async function er() {
	if (Ae) return new Promise((e) => {
		requestAnimationFrame(() => e()), setTimeout(() => e());
	});
	await Promise.resolve(), at();
}
function q(e) {
	var t = (e.f & 2) != 0;
	if (jn?.add(e), W !== null && !Fn && !(G !== null && G.f & 16384) && (Rn === null || !n.call(Rn, e))) {
		var r = W.deps;
		if (W.f & 2097152) e.rv < Wn && (e.rv = Wn, K === null && r !== null && r[Bn] === e ? Bn++ : K === null ? K = [e] : K.push(e));
		else {
			(W.deps ??= []).push(e);
			var i = e.reactions;
			i === null ? e.reactions = [W] : n.call(i, W) || i.push(W);
		}
	}
	if (Nn && Mt.has(e)) return Mt.get(e);
	if (t) {
		var a = e;
		if (Nn) {
			var o = a.v;
			return (!(a.f & 1024) && a.reactions !== null || nr(a)) && (o = Dt(a)), Mt.set(a, o), o;
		}
		var s = (a.f & 512) == 0 && !Fn && W !== null && (Mn || (W.f & 512) != 0), c = (a.f & y) === 0;
		Jn(a) && (s && (a.f |= 512), Ot(a)), s && !c && (At(a), tr(a));
	}
	if (I?.has(e)) return I.get(e);
	if (e.f & 8388608) throw e.v;
	return e.v;
}
function tr(e) {
	if (e.f |= 512, e.deps !== null) for (let t of e.deps) (t.reactions ??= []).push(e), t.f & 2 && !(t.f & 512) && (At(t), tr(t));
}
function nr(e) {
	if (e.v === k) return !0;
	if (e.deps === null) return !1;
	for (let t of e.deps) if (Mt.has(t) || t.f & 2 && nr(t)) return !0;
	return !1;
}
function rr(e) {
	var t = Fn;
	try {
		return Fn = !0, e();
	} finally {
		Fn = t;
	}
}
[.../* @__PURE__ */ "allowfullscreen.async.autofocus.autoplay.checked.controls.default.disabled.formnovalidate.indeterminate.inert.ismap.loop.multiple.muted.nomodule.novalidate.open.playsinline.readonly.required.reversed.seamless.selected.webkitdirectory.defer.disablepictureinpicture.disableremoteplayback".split(".")];
var ir = ["touchstart", "touchmove"];
function ar(e) {
	return ir.includes(e);
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/events.js
var or = Symbol("events"), sr = /* @__PURE__ */ new Set(), cr = /* @__PURE__ */ new Set();
function lr(e, t, n, r = {}) {
	function i(e) {
		if (r.capture || pr.call(t, e), !e.cancelBubble) return on(() => n?.call(this, e));
	}
	return e.startsWith("pointer") || e.startsWith("touch") || e === "wheel" ? Re(() => {
		t.addEventListener(e, i, r);
	}) : t.addEventListener(e, i, r), i;
}
function ur(e, t, n, r, i) {
	var a = {
		capture: r,
		passive: i
	}, o = lr(e, t, n, a);
	(t === document.body || t === window || t === document || t instanceof HTMLMediaElement) && fn(() => {
		t.removeEventListener(e, o, a);
	});
}
function J(e, t, n) {
	(t[or] ??= {})[e] = n;
}
function dr(e) {
	for (var t = 0; t < e.length; t++) sr.add(e[t]);
	for (var n of cr) n(e);
}
var fr = null;
function pr(e) {
	var t = this, n = t.ownerDocument, r = e.type, a = e.composedPath?.() || [], o = a[0] || e.target;
	fr = e;
	var s = 0, c = fr === e && e[or];
	if (c) {
		var l = a.indexOf(c);
		if (l !== -1 && (t === document || t === window)) {
			e[or] = t;
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
		var d = W, f = G;
		In(null), Ln(null);
		try {
			for (var p, m = []; o !== null;) {
				var h = o.assignedSlot || o.parentNode || o.host || null;
				try {
					var g = o[or]?.[r];
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
			e[or] = t, delete e.currentTarget, In(d), Ln(f);
		}
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/reconciler.js
var mr = globalThis?.window?.trustedTypes && /* @__PURE__ */ globalThis.window.trustedTypes.createPolicy("svelte-trusted-html", { createHTML: (e) => e });
function hr(e) {
	return mr?.createHTML(e) ?? e;
}
function gr(e) {
	var t = en("template");
	return t.innerHTML = hr(e.replaceAll("<!>", "<!---->")), t.content;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/template.js
function _r(e, t) {
	var n = G;
	n.nodes === null && (n.nodes = {
		start: e,
		end: t,
		a: null,
		t: null
	});
}
/* @__NO_SIDE_EFFECTS__ */
function Y(e, t) {
	var n = (t & 1) != 0, r = (t & 2) != 0, i, a = !e.startsWith("<!>");
	return () => {
		if (A) return _r(j, null), j;
		i === void 0 && (i = gr(a ? e : "<!>" + e), n || (i = /* @__PURE__ */ Yt(i)));
		var t = r || Wt ? document.importNode(i, !0) : i.cloneNode(!0);
		if (n) {
			var o = /* @__PURE__ */ Yt(t), s = t.lastChild;
			_r(o, s);
		} else _r(t, t);
		return t;
	};
}
function vr() {
	if (A) return _r(j, null), j;
	var e = document.createDocumentFragment(), t = document.createComment(""), n = Jt();
	return e.append(t, n), _r(t, n), e;
}
function X(e, t) {
	if (A) {
		var n = G;
		(!(n.f & 32768) || n.nodes.end === null) && (n.nodes.end = j), Ce();
		return;
	}
	e !== null && e.before(t);
}
function Z(e, t) {
	var n = t == null ? "" : typeof t == "object" ? `${t}` : t;
	n !== (e.__t ??= e.nodeValue) && (e.__t = n, e.nodeValue = `${n}`);
}
function yr(e, t) {
	return xr(e, t);
}
var br = /* @__PURE__ */ new Map();
function xr(e, { target: t, anchor: n, props: i = {}, events: a, context: o, intro: s = !0, transformError: c }) {
	qt();
	var l = void 0, u = hn(() => {
		var s = n ?? t.appendChild(Jt());
		gt(s, { pending: () => {} }, (t) => {
			Ne({});
			var n = N;
			if (o && (n.c = o), a && (i.$$events = a), A && _r(t, null), l = e(t, i) || {}, A && (G.nodes.end = j, j === null || j.nodeType !== 8 || j.data !== "]")) throw ve(), ge;
			Pe();
		}, c);
		var u = /* @__PURE__ */ new Set(), d = (e) => {
			for (var n = 0; n < e.length; n++) {
				var r = e[n];
				if (!u.has(r)) {
					u.add(r);
					var i = ar(r);
					for (let e of [t, document]) {
						var a = br.get(e);
						a === void 0 && (a = /* @__PURE__ */ new Map(), br.set(e, a));
						var o = a.get(r);
						o === void 0 ? (e.addEventListener(r, pr, { passive: i }), a.set(r, 1)) : a.set(r, o + 1);
					}
				}
			}
		};
		return d(r(sr)), cr.add(d), () => {
			for (var e of u) for (let n of [t, document]) {
				var r = br.get(n), i = r.get(e);
				--i == 0 ? (n.removeEventListener(e, pr), r.delete(e), r.size === 0 && br.delete(n)) : r.set(e, i);
			}
			cr.delete(d), s !== n && s.parentNode?.removeChild(s);
		};
	});
	return Sr.set(l, u), l;
}
var Sr = /* @__PURE__ */ new WeakMap(), Cr = class {
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
			if (n) On(n), this.#r.delete(t);
			else {
				var r = this.#n.get(t);
				r && (this.#t.set(t, r.effect), this.#n.delete(t), r.fragment.lastChild.remove(), this.anchor.before(r.fragment), n = r.effect);
			}
			for (let [t, n] of this.#e) {
				if (this.#e.delete(t), t === e) break;
				let r = this.#n.get(n);
				r && (U(r.effect), this.#n.delete(n));
			}
			for (let [e, r] of this.#t) {
				if (e === t || this.#r.has(e)) continue;
				let i = () => {
					if (Array.from(this.#e.values()).includes(e)) {
						var t = document.createDocumentFragment();
						An(r, t), t.append(Jt()), this.#n.set(e, {
							effect: r,
							fragment: t
						});
					} else U(r);
					this.#r.delete(e), this.#t.delete(e);
				};
				this.#i || !n ? (this.#r.add(e), En(r, i, !1)) : i();
			}
		}
	};
	#o = (e) => {
		this.#e.delete(e);
		let t = Array.from(this.#e.values());
		for (let [e, n] of this.#n) t.includes(e) || (U(n.effect), this.#n.delete(e));
	};
	ensure(e, t) {
		var n = F, r = $t();
		if (t && !this.#t.has(e) && !this.#n.has(e)) if (r) {
			var i = document.createDocumentFragment(), a = Jt();
			i.append(a), this.#n.set(e, {
				effect: bn(() => t(a)),
				fragment: i
			});
		} else this.#t.set(e, bn(() => t(this.anchor)));
		if (this.#e.set(n, e), r) {
			for (let [t, r] of this.#t) t === e ? n.unskip_effect(r) : n.skip_effect(r);
			for (let [t, r] of this.#n) t === e ? n.unskip_effect(r.effect) : n.skip_effect(r.effect);
			n.oncommit(this.#a), n.ondiscard(this.#o);
		} else A && (this.anchor = j), this.#a(n);
	}
};
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/if.js
function Q(e, t, n = !1) {
	var r;
	A && (r = j, Ce());
	var i = new Cr(e), a = n ? x : 0;
	function o(e, t) {
		if (A) {
			var n = Ee(r);
			if (e !== parseInt(n.substring(1))) {
				var a = Te();
				Se(a), i.anchor = a, xe(!1), i.ensure(e, t), xe(!0);
				return;
			}
		}
		i.ensure(e, t);
	}
	yn(() => {
		var e = !1;
		t((t, n = 0) => {
			e = !0, o(n, t);
		}), e || o(-1, null);
	}, a);
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/each.js
function wr(e, t) {
	return t;
}
function Tr(e, t, n) {
	for (var i = [], a = t.length, o, s = t.length, c = 0; c < a; c++) {
		let n = t[c];
		En(n, () => {
			if (o) {
				if (o.pending.delete(n), o.done.add(n), o.pending.size === 0) {
					var t = e.outrogroups;
					Er(e, r(o.done)), t.delete(o), t.size === 0 && (e.outrogroups = null);
				}
			} else --s;
		}, !1);
	}
	if (s === 0) {
		var l = i.length === 0 && n !== null;
		if (l) {
			var u = n, d = u.parentNode;
			Qt(d), d.append(u), e.items.clear();
		}
		Er(e, t, !l);
	} else o = {
		pending: new Set(t),
		done: /* @__PURE__ */ new Set()
	}, (e.outrogroups ??= /* @__PURE__ */ new Set()).add(o);
}
function Er(e, t, n = !0) {
	var r;
	if (e.pending.size > 0) {
		r = /* @__PURE__ */ new Set();
		for (let t of e.pending.values()) for (let n of t) r.add(e.items.get(n).e);
	}
	for (var i = 0; i < t.length; i++) {
		var a = t[i];
		r?.has(a) ? (a.f |= C, An(a, document.createDocumentFragment())) : U(t[i], n);
	}
}
var Dr;
function Or(t, n, i, a, o, s = null) {
	var c = t, l = /* @__PURE__ */ new Map();
	if (n & 4) {
		var u = t;
		c = A ? Se(/* @__PURE__ */ Yt(u)) : u.appendChild(Jt());
	}
	A && Ce();
	var d = null, f = /* @__PURE__ */ wt(() => {
		var t = i();
		return e(t) ? t : t == null ? [] : r(t);
	}), p, m = /* @__PURE__ */ new Map(), h = !0;
	function g(e) {
		v.effect.f & 16384 || (v.pending.delete(e), v.fallback = d, Ar(v, p, c, n, a), d !== null && (p.length === 0 ? d.f & 33554432 ? (d.f ^= C, Mr(d, null, c)) : On(d) : En(d, () => {
			d = null;
		})));
	}
	function _(e) {
		v.pending.delete(e);
	}
	var v = {
		effect: yn(() => {
			p = q(f);
			var e = p.length;
			let t = !1;
			A && Ee(c) === "[!" != (e === 0) && (c = Te(), Se(c), xe(!1), t = !0);
			for (var r = /* @__PURE__ */ new Set(), u = F, v = $t(), y = 0; y < e; y += 1) {
				A && j.nodeType === 8 && j.data === "]" && (c = j, t = !0, xe(!1));
				var b = p[y], x = a(b, y), S = h ? null : l.get(x);
				S ? (S.v && It(S.v, b), S.i && It(S.i, y), v && u.unskip_effect(S.e)) : (S = jr(l, h ? c : Dr ??= Jt(), b, x, y, o, n, i), h || (S.e.f |= C), l.set(x, S)), r.add(x);
			}
			if (e === 0 && s && !d && (h ? d = bn(() => s(c)) : (d = bn(() => s(Dr ??= Jt())), d.f |= C)), e > r.size && oe("", "", ""), A && e > 0 && Se(Te()), !h) if (m.set(u, r), v) {
				for (let [e, t] of l) r.has(e) || u.skip_effect(t.e);
				u.oncommit(g), u.ondiscard(_);
			} else g(u);
			t && xe(!0), q(f);
		}),
		flags: n,
		items: l,
		pending: m,
		outrogroups: null,
		fallback: d
	};
	h = !1, A && (c = j);
}
function kr(e) {
	for (; e !== null && !(e.f & 32);) e = e.next;
	return e;
}
function Ar(e, t, n, i, a) {
	var o = (i & 8) != 0, s = t.length, c = e.items, l = kr(e.effect.first), u, d = null, f, p = [], m = [], h, g, _, v;
	if (o) for (v = 0; v < s; v += 1) h = t[v], g = a(h, v), _ = c.get(g).e, _.f & 33554432 || (_.nodes?.a?.measure(), (f ??= /* @__PURE__ */ new Set()).add(_));
	for (v = 0; v < s; v += 1) {
		if (h = t[v], g = a(h, v), _ = c.get(g).e, e.outrogroups !== null) for (let t of e.outrogroups) t.pending.delete(_), t.done.delete(_);
		if (_.f & 33554432) if (_.f ^= C, _ === l) Mr(_, null, n);
		else {
			var y = d ? d.next : l;
			_ === e.effect.last && (e.effect.last = _.prev), _.prev && (_.prev.next = _.next), _.next && (_.next.prev = _.prev), Nr(e, d, _), Nr(e, _, y), Mr(_, y, n), d = _, p = [], m = [], l = kr(d.next);
			continue;
		}
		if (_.f & 8192 && (On(_), o && (_.nodes?.a?.unfix(), (f ??= /* @__PURE__ */ new Set()).delete(_))), _ !== l) {
			if (u !== void 0 && u.has(_)) {
				if (p.length < m.length) {
					var b = m[0], x;
					d = b.prev;
					var S = p[0], ee = p[p.length - 1];
					for (x = 0; x < p.length; x += 1) Mr(p[x], b, n);
					for (x = 0; x < m.length; x += 1) u.delete(m[x]);
					Nr(e, S.prev, ee.next), Nr(e, d, S), Nr(e, ee, b), l = b, d = ee, --v, p = [], m = [];
				} else u.delete(_), Mr(_, l, n), Nr(e, _.prev, _.next), Nr(e, _, d === null ? e.effect.first : d.next), Nr(e, d, _), d = _;
				continue;
			}
			for (p = [], m = []; l !== null && l !== _;) (u ??= /* @__PURE__ */ new Set()).add(l), m.push(l), l = kr(l.next);
			if (l === null) continue;
		}
		_.f & 33554432 || p.push(_), d = _, l = kr(_.next);
	}
	if (e.outrogroups !== null) {
		for (let t of e.outrogroups) t.pending.size === 0 && (Er(e, r(t.done)), e.outrogroups?.delete(t));
		e.outrogroups.size === 0 && (e.outrogroups = null);
	}
	if (l !== null || u !== void 0) {
		var w = [];
		if (u !== void 0) for (_ of u) _.f & 8192 || w.push(_);
		for (; l !== null;) !(l.f & 8192) && l !== e.fallback && w.push(l), l = kr(l.next);
		var te = w.length;
		if (te > 0) {
			var T = i & 4 && s === 0 ? n : null;
			if (o) {
				for (v = 0; v < te; v += 1) w[v].nodes?.a?.measure();
				for (v = 0; v < te; v += 1) w[v].nodes?.a?.fix();
			}
			Tr(e, w, T);
		}
	}
	o && Re(() => {
		if (f !== void 0) for (_ of f) _.nodes?.a?.apply();
	});
}
function jr(e, t, n, r, i, a, o, s) {
	var c = o & 1 ? o & 16 ? Pt(n) : /* @__PURE__ */ Ft(n, !1, !1) : null, l = o & 2 ? Pt(i) : null;
	return {
		v: c,
		i: l,
		e: bn(() => (a(t, c ?? n, l ?? i, s), () => {
			e.delete(r);
		}))
	};
}
function Mr(e, t, n) {
	if (e.nodes) for (var r = e.nodes.start, i = e.nodes.end, a = t && !(t.f & 33554432) ? t.nodes.start : n; r !== null;) {
		var o = /* @__PURE__ */ Xt(r);
		if (a.before(r), r === i) return;
		r = o;
	}
}
function Nr(e, t, n) {
	t === null ? e.effect.first = n : t.next = n, n === null ? e.effect.last = t : n.prev = t;
}
//#endregion
//#region node_modules/svelte/src/internal/shared/attributes.js
var Pr = [..." 	\n\r\f\xA0\v﻿"];
function Fr(e, t, n) {
	var r = e == null ? "" : "" + e;
	if (t && (r = r ? r + " " + t : t), n) {
		for (var i of Object.keys(n)) if (n[i]) r = r ? r + " " + i : i;
		else if (r.length) for (var a = i.length, o = 0; (o = r.indexOf(i, o)) >= 0;) {
			var s = o + a;
			(o === 0 || Pr.includes(r[o - 1])) && (s === r.length || Pr.includes(r[s])) ? r = (o === 0 ? "" : r.substring(0, o)) + r.substring(s + 1) : o = s;
		}
	}
	return r === "" ? null : r;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/class.js
function Ir(e, t, n, r, i, a) {
	var o = e.__className;
	if (A || o !== n || o === void 0) {
		var s = Fr(n, r, a);
		(!A || s !== e.getAttribute("class")) && (s == null ? e.removeAttribute("class") : t ? e.className = s : e.setAttribute("class", s)), e.__className = n;
	} else if (a && i !== a) for (var c in a) {
		var l = !!a[c];
		(i == null || l !== !!i[c]) && e.classList.toggle(c, l);
	}
	return a;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/select.js
function Lr(t, n, r = !1) {
	if (t.multiple) {
		if (n == null) return;
		if (!e(n)) return ye();
		for (var i of t.options) i.selected = n.includes(zr(i));
		return;
	}
	for (i of t.options) if (Ht(zr(i), n)) {
		i.selected = !0;
		return;
	}
	(!r || n !== void 0) && (t.selectedIndex = -1);
}
function Rr(e) {
	var t = new MutationObserver(() => {
		Lr(e, e.__value);
	});
	t.observe(e, {
		childList: !0,
		subtree: !0,
		attributes: !0,
		attributeFilter: ["value"]
	}), fn(() => {
		t.disconnect();
	});
}
function zr(e) {
	return "__value" in e ? e.__value : e.value;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/attributes.js
var Br = Symbol("is custom element"), Vr = Symbol("is html"), Hr = ie ? "link" : "LINK", Ur = ie ? "progress" : "PROGRESS";
function $(e) {
	if (A) {
		var t = !1, n = () => {
			if (!t) {
				if (t = !0, e.hasAttribute("value")) {
					var n = e.value;
					Kr(e, "value", null), e.value = n;
				}
				if (e.hasAttribute("checked")) {
					var r = e.checked;
					Kr(e, "checked", null), e.checked = r;
				}
			}
		};
		e.__on_r = n, Re(n), an();
	}
}
function Wr(e, t) {
	var n = qr(e);
	n.value === (n.value = t ?? void 0) || e.value === t && (t !== 0 || e.nodeName !== Ur) || (e.value = t ?? "");
}
function Gr(e, t) {
	var n = qr(e);
	n.checked !== (n.checked = t ?? void 0) && (e.checked = t);
}
function Kr(e, t, n, r) {
	var i = qr(e);
	A && (i[t] = e.getAttribute(t), t === "src" || t === "srcset" || t === "href" && e.nodeName === Hr) || i[t] !== (i[t] = n) && (t === "loading" && (e[re] = n), n == null ? e.removeAttribute(t) : typeof n != "string" && Yr(e).includes(t) ? e[t] = n : e.setAttribute(t, n));
}
function qr(e) {
	return e.__attributes ??= {
		[Br]: e.nodeName.includes("-"),
		[Vr]: e.namespaceURI === _e
	};
}
var Jr = /* @__PURE__ */ new Map();
function Yr(e) {
	var t = e.getAttribute("is") || e.nodeName, n = Jr.get(t);
	if (n) return n;
	Jr.set(t, n = []);
	for (var r, i = e, a = Element.prototype; a !== i;) {
		for (var s in r = o(i), r) r[s].set && n.push(s);
		i = l(i);
	}
	return n;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/input.js
function Xr(e, t, n = t) {
	var r = /* @__PURE__ */ new WeakSet();
	sn(e, "input", async (i) => {
		var a = i ? e.defaultValue : e.value;
		if (a = Zr(e) ? Qr(a) : a, n(a), F !== null && r.add(F), await er(), a !== (a = t())) {
			var o = e.selectionStart, s = e.selectionEnd, c = e.value.length;
			if (e.value = a ?? "", s !== null) {
				var l = e.value.length;
				o === s && s === c && l > c ? (e.selectionStart = l, e.selectionEnd = l) : (e.selectionStart = o, e.selectionEnd = Math.min(s, l));
			}
		}
	}), (A && e.defaultValue !== e.value || rr(t) == null && e.value) && (n(Zr(e) ? Qr(e.value) : e.value), F !== null && r.add(F)), vn(() => {
		var n = t();
		if (e === document.activeElement) {
			var i = Ae ? Xe : F;
			if (r.has(i)) return;
		}
		Zr(e) && n === Qr(e.value) || e.type === "date" && !n && !e.value || n !== e.value && (e.value = n ?? "");
	});
}
function Zr(e) {
	var t = e.type;
	return t === "number" || t === "range";
}
function Qr(e) {
	return e === "" ? null : +e;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/this.js
function $r(e, t) {
	return e === t || e?.[D] === t;
}
function ei(e = {}, t, n, r) {
	var i = N.r, a = G;
	return gn(() => {
		var o, s;
		return vn(() => {
			o = s, s = r?.() || [], rr(() => {
				e !== n(...s) && (t(e, ...s), o && $r(n(...o), e) && t(null, ...o));
			});
		}), () => {
			let r = a;
			for (; r !== i && r.parent !== null && r.parent.f & 33554432;) r = r.parent;
			let o = () => {
				s && $r(n(...s), e) && t(null, ...s);
			}, c = r.teardown;
			r.teardown = () => {
				o(), c?.();
			};
		};
	}), e;
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/props.js
function ti(e, t, n, r) {
	var i = !je || (n & 2) != 0, o = (n & 8) != 0, s = (n & 16) != 0, c = r, l = !0, u = () => (l && (l = !1, c = s ? rr(r) : r), c);
	let d;
	if (o) {
		var f = D in e || ne in e;
		d = a(e, t)?.set ?? (f && t in e ? (n) => e[t] = n : void 0);
	}
	var p, m = !1;
	o ? [p, m] = Je(() => e[t]) : p = e[t], p === void 0 && r !== void 0 && (p = u(), d && (i && de(t), d(p)));
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
	var _ = !1, v = (n & 1 ? St : wt)(() => (_ = !1, h()));
	o && q(v);
	var y = G;
	return (function(e, t) {
		if (arguments.length > 0) {
			let n = t ? q(v) : i && o ? Bt(e) : e;
			return z(v, n), _ = !0, c !== void 0 && (c = n), e;
		}
		return Nn && _ || y.f & 16384 ? v.v : q(v);
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
var ui = /* @__PURE__ */ Y("<option> </option>"), di = /* @__PURE__ */ Y("<tr class=\"border-b border-white/[0.06]\"><td colspan=\"9\" class=\"px-2 py-1.5\"><div class=\"flex items-center gap-2 pl-6\"><svg class=\"w-3 h-3 text-white/20 shrink-0\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M3 10h10a3 3 0 013 3v1\"></path></svg> <input type=\"text\" placeholder=\"Add a description or note...\" class=\"w-full px-2 py-1 text-xs bg-transparent border-0 border-b border-white/[0.06] text-[var(--color-concrete)]\n						focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/20 font-[var(--font-body)]\"/></div></td></tr>"), fi = /* @__PURE__ */ Y("<tr class=\"border-b border-white/[0.06] hover:bg-white/[0.03] text-sm group\"><td class=\"px-1 py-1 w-24\"><select></select></td><td class=\"px-1 py-1\"><input type=\"text\" class=\"w-full px-1 py-0.5 text-[var(--color-white)] bg-transparent border-0 rounded\n				focus:ring-2 focus:ring-[var(--color-sunburst)] font-[var(--font-body)]\" placeholder=\"Item name\"/></td><td class=\"px-1 py-1 w-20\"><input type=\"number\" step=\"any\" min=\"0\" class=\"w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 rounded text-[var(--color-white)]\n				focus:ring-2 focus:ring-[var(--color-sunburst)]\"/></td><td class=\"px-1 py-1 w-16\"><input type=\"text\" class=\"w-full text-center text-[var(--color-muted-text)] px-1 py-0.5 bg-transparent border-0 rounded\n				focus:ring-2 focus:ring-[var(--color-sunburst)] text-sm font-[var(--font-body)]\" placeholder=\"ea\"/></td><td class=\"px-1 py-1 w-24\"><input type=\"number\" step=\"0.01\" min=\"0\" class=\"w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 rounded text-[var(--color-white)]\n				focus:ring-2 focus:ring-[var(--color-sunburst)]\"/></td><td class=\"px-2 py-1.5 text-right font-mono text-white/40 w-16 text-xs\"> </td><td class=\"px-2 py-1.5 text-right font-mono text-[var(--color-muted-text)] w-24 text-xs\"> </td><td class=\"px-2 py-1.5 text-right font-mono font-medium text-[var(--color-white)] w-24\"> </td><td class=\"px-1 py-1 w-16\"><div class=\"flex items-center gap-0.5\"><button title=\"Toggle description\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z\"></path></svg></button> <button class=\"opacity-0 group-hover:opacity-100 text-white/30 hover:text-red-400 transition-opacity p-0.5\" title=\"Delete item\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></button></div></td></tr> <!>", 1);
function pi(e, t) {
	Ne(t, !0);
	let n = ti(t, "item", 7), r = [
		"materials",
		"labor",
		"equipment",
		"subs",
		"other"
	], i = {
		materials: "bg-blue-900/30 text-blue-400",
		labor: "bg-amber-900/30 text-amber-400",
		equipment: "bg-purple-900/30 text-purple-400",
		subs: "bg-green-900/30 text-green-400",
		other: "bg-white/[0.06] text-[var(--color-concrete)]"
	}, a = /* @__PURE__ */ L(() => ni(n().category_type, t.globals, t.markupOverrides, t.markupEnabled)), o = /* @__PURE__ */ L(() => ri(n().unit_price, q(a))), s = /* @__PURE__ */ L(() => ii(n().quantity, n().unit_price, q(a))), c = /* @__PURE__ */ L(() => i[n().category_type] || i.other), l = /* @__PURE__ */ R(!!n().description);
	function u() {
		t.onchange?.();
	}
	function d(e) {
		let t = parseFloat(e.target.value);
		isNaN(t) || (n().quantity = t, u());
	}
	function f(e) {
		let t = parseFloat(e.target.value);
		isNaN(t) || (n().unit_price = t, n().price_override = !0, u());
	}
	function p(e) {
		n().item_name = e.target.value, u();
	}
	function m(e) {
		n().unit = e.target.value, u();
	}
	function h(e) {
		n().category_type = e.target.value, u();
	}
	function g(e) {
		n().description = e.target.value || null, u();
	}
	function _() {
		z(l, !q(l)), q(l);
	}
	function v() {
		t.ondelete?.(n().id);
	}
	var y = fi(), b = Zt(y), x = B(b), S = B(x);
	Or(S, 21, () => r, wr, (e, t) => {
		var n = ui(), r = B(n, !0);
		M(n);
		var i = {};
		H(() => {
			Z(r, q(t)), i !== (i = q(t)) && (n.value = (n.__value = q(t)) ?? "");
		}), X(e, n);
	}), M(S);
	var ee;
	Rr(S), M(x);
	var C = V(x), w = B(C);
	$(w), M(C);
	var te = V(C), T = B(te);
	$(T), M(te);
	var E = V(te), D = B(E);
	$(D), M(E);
	var ne = V(E), re = B(ne);
	$(re), M(ne);
	var O = V(ne), ie = B(O);
	M(O);
	var ae = V(O), oe = B(ae, !0);
	M(ae);
	var se = V(ae), ce = B(se, !0);
	M(se);
	var le = V(se), ue = B(le), de = B(ue), fe = V(de, 2);
	M(ue), M(le), M(b);
	var pe = V(b, 2), me = (e) => {
		var t = di(), r = B(t), i = B(r), a = V(B(i), 2);
		$(a), M(i), M(r), M(t), H(() => Wr(a, n().description ?? "")), J("input", a, g), X(e, t);
	};
	Q(pe, (e) => {
		q(l) && e(me);
	}), H((e, t) => {
		Ir(S, 1, `w-full text-xs px-1 py-1 rounded border-0 bg-transparent font-medium cursor-pointer
				focus:ring-2 focus:ring-[var(--color-sunburst)] ${q(c) ?? ""} font-[var(--font-ui)]`), ee !== (ee = n().category_type) && (S.value = (S.__value = n().category_type) ?? "", Lr(S, n().category_type)), Wr(w, n().item_name), Wr(T, n().quantity), Wr(D, n().unit), Wr(re, n().unit_price), Z(ie, `${q(a) ?? ""}%`), Z(oe, e), Z(ce, t), Ir(de, 1, `opacity-0 group-hover:opacity-100 p-0.5 transition-opacity rounded
					${q(l) || n().description ? "opacity-100 text-[var(--color-sunburst)]" : "text-white/30 hover:text-[var(--color-concrete)]"}`);
	}, [() => ci(q(o)), () => ci(q(s))]), J("change", S, h), J("input", w, p), J("input", T, d), J("input", D, m), J("input", re, f), J("click", de, _), J("click", fe, v), X(e, y), Pe();
}
dr([
	"change",
	"input",
	"click"
]);
//#endregion
//#region src/lib/Autocomplete.svelte
var mi = /* @__PURE__ */ Y("<button><div><span class=\"text-[var(--color-white)] font-[var(--font-body)]\"> </span> <span class=\"text-xs text-white/40 ml-2 font-[var(--font-body)]\"> </span></div> <div class=\"flex items-center gap-2 text-xs\"><span class=\"font-mono text-[var(--color-muted-text)]\"> </span> <span class=\"text-white/40\"> </span></div></button>"), hi = /* @__PURE__ */ Y("<div class=\"absolute top-full left-0 right-0 mt-1 bg-[var(--color-granite)] border border-white/[0.08] rounded-lg shadow-lg z-50 max-h-64 overflow-y-auto\"></div>"), gi = /* @__PURE__ */ Y("<div class=\"absolute top-full left-0 right-0 mt-1 bg-[var(--color-granite)] border border-white/[0.08] rounded-lg shadow-lg z-50\"><button class=\"w-full text-left px-3 py-2 text-sm text-[var(--color-muted-text)] hover:bg-white/[0.06] font-[var(--font-body)]\">Add custom item: <span class=\"font-medium text-[var(--color-white)]\"> </span></button></div>"), _i = /* @__PURE__ */ Y("<div class=\"relative\"><input type=\"text\" placeholder=\"Search items or type a name...\" class=\"w-full px-3 py-2 text-sm bg-white/[0.04] border border-white/[0.08] rounded-lg text-[var(--color-white)]\n			focus:ring-2 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] outline-none placeholder-white/30 font-[var(--font-body)]\"/> <!></div>");
function vi(e, t) {
	Ne(t, !0);
	let n = ti(t, "materialsDb", 19, () => []), r = ti(t, "ratesDb", 19, () => []), i = ti(t, "categoryType", 3, "materials"), a = /* @__PURE__ */ R(""), o = /* @__PURE__ */ R(0), s, c = /* @__PURE__ */ L(() => {
		if (q(a).length < 1) return [];
		let e = q(a).toLowerCase(), t;
		return t = i() === "materials" ? n().map((e) => ({
			id: e.id,
			name: e.name,
			category: e.supplier || "",
			unit: e.unit,
			price: e.unit_price,
			type: "materials",
			source: "material"
		})) : r().filter((e) => i() === "labor" ? e.category === "Labor" : i() === "equipment" ? e.category === "Equipment Rentals" : i() === "subs" ? e.category === "Subcontractors" : i() === "other" ? e.category === "Other" : !0).map((e) => ({
			id: e.id,
			name: e.name,
			category: e.category,
			unit: e.unit,
			price: e.rate,
			type: i(),
			source: "rate"
		})), t.filter((t) => t.name.toLowerCase().includes(e) || t.category.toLowerCase().includes(e)).slice(0, 15);
	});
	function l(e) {
		e.key === "ArrowDown" ? (e.preventDefault(), z(o, Math.min(q(o) + 1, q(c).length - 1), !0)) : e.key === "ArrowUp" ? (e.preventDefault(), z(o, Math.max(q(o) - 1, 0), !0)) : e.key === "Enter" ? (e.preventDefault(), q(c).length > 0 && q(o) < q(c).length ? u(q(c)[q(o)]) : q(a).trim() && d()) : e.key === "Escape" && (e.preventDefault(), t.oncancel?.());
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
			item_name: q(a).trim(),
			unit: "ea",
			unit_price: 0,
			category_type: i(),
			is_custom: !0,
			material_id: null
		});
	}
	pn(() => {
		z(o, 0);
	}), pn(() => {
		s?.focus();
	});
	var f = _i(), p = B(f);
	$(p), ei(p, (e) => s = e, () => s);
	var m = V(p, 2), h = (e) => {
		var t = hi();
		Or(t, 21, () => q(c), wr, (e, t, n) => {
			var r = mi(), i = B(r), a = B(i), s = B(a, !0);
			M(a);
			var c = V(a, 2), l = B(c, !0);
			M(c), M(i);
			var d = V(i, 2), f = B(d), p = B(f);
			M(f);
			var m = V(f, 2), h = B(m);
			M(m), M(d), M(r), H((e) => {
				Ir(r, 1, `w-full text-left px-3 py-2 text-sm flex items-center justify-between
						hover:bg-white/[0.06] ${n === q(o) ? "bg-[var(--color-sunburst)]/10" : ""}`), Z(s, q(t).name), Z(l, q(t).category), Z(p, `$${e ?? ""}`), Z(h, `/ ${q(t).unit ?? ""}`);
			}, [() => q(t).price.toFixed(2)]), J("click", r, () => u(q(t))), ur("mouseenter", r, () => z(o, n, !0)), X(e, r);
		}), M(t), X(e, t);
	}, g = (e) => {
		var t = gi(), n = B(t), r = V(B(n)), i = B(r);
		M(r), M(n), M(t), H(() => Z(i, `"${q(a) ?? ""}"`)), J("click", n, d), X(e, t);
	};
	Q(m, (e) => {
		q(c).length > 0 ? e(h) : q(a).length > 0 && e(g, 1);
	}), M(f), J("keydown", p, l), Xr(p, () => q(a), (e) => z(a, e)), X(e, f), Pe();
}
dr(["keydown", "click"]);
//#endregion
//#region node_modules/nanoid/url-alphabet/index.js
var yi = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict", bi = (e = 21) => {
	let t = "", n = crypto.getRandomValues(new Uint8Array(e |= 0));
	for (; e--;) t += yi[n[e] & 63];
	return t;
}, xi = /* @__PURE__ */ Y("<input type=\"text\" class=\"text-xs font-semibold text-[var(--color-concrete)] uppercase tracking-wide px-1 py-0.5 border border-white/[0.08] rounded bg-white/[0.04] focus:ring-2 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] font-[var(--font-ui)]\"/>"), Si = /* @__PURE__ */ Y("<span role=\"button\" tabindex=\"0\" class=\"text-xs font-semibold text-[var(--color-muted-text)] uppercase tracking-wide cursor-pointer font-[var(--font-ui)]\"> </span>"), Ci = /* @__PURE__ */ Y("<div class=\"mb-2\"><!></div>"), wi = /* @__PURE__ */ Y("<table class=\"w-full\"><tbody></tbody></table>"), Ti = /* @__PURE__ */ Y("<div class=\"ml-4 mt-2\"><div class=\"flex items-center gap-2 mb-1 group/cg\"><!> <span class=\"text-xs text-white/40 font-[var(--font-body)]\"> </span> <button class=\"opacity-0 group-hover/cg:opacity-100 text-white/30 hover:text-red-400 transition-opacity p-0.5\" title=\"Delete group\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></button> <button class=\"text-xs text-[var(--color-sunburst)] hover:brightness-110 ml-auto font-[var(--font-ui)]\">+ Add Item</button></div> <!> <!></div>");
function Ei(e, t) {
	Ne(t, !0);
	let n = ti(t, "group", 7), r = /* @__PURE__ */ R(!1), i = /* @__PURE__ */ R(!1), a = /* @__PURE__ */ R(Bt(n().name));
	function o(e) {
		t.onsnapshot?.(), n().line_items.push({
			id: bi(),
			category_type: e.category_type,
			item_name: e.item_name,
			quantity: 1,
			unit: e.unit,
			unit_price: e.unit_price,
			is_custom: e.is_custom,
			material_id: e.material_id,
			price_override: !1,
			description: null,
			sort_order: n().line_items.length,
			component_group_id: n().id
		}), z(r, !1), t.onchange?.();
	}
	function s(e) {
		t.onsnapshot?.();
		let r = n().line_items.findIndex((t) => t.id === e);
		r !== -1 && (n().line_items.splice(r, 1), t.onchange?.());
	}
	function c() {
		z(a, n().name, !0), z(i, !0);
	}
	function l() {
		let e = q(a).trim();
		e && e !== n().name && (t.onsnapshot?.(), n().name = e, t.onchange?.()), z(i, !1);
	}
	var u = Ti(), d = B(u), f = B(d), p = (e) => {
		var t = xi();
		$(t), nn(t, !0), J("keydown", t, (e) => {
			e.key === "Enter" && l(), e.key === "Escape" && z(i, !1);
		}), ur("blur", t, l), Xr(t, () => q(a), (e) => z(a, e)), X(e, t);
	}, m = (e) => {
		var t = Si(), r = B(t, !0);
		M(t), H(() => Z(r, n().name)), J("dblclick", t, c), X(e, t);
	};
	Q(f, (e) => {
		q(i) ? e(p) : e(m, -1);
	});
	var h = V(f, 2), g = B(h);
	M(h);
	var _ = V(h, 2), v = V(_, 2);
	M(d);
	var y = V(d, 2), b = (e) => {
		var n = Ci();
		vi(B(n), {
			get materialsDb() {
				return t.materialsDb;
			},
			get ratesDb() {
				return t.ratesDb;
			},
			categoryType: "materials",
			onselect: o,
			oncancel: () => z(r, !1)
		}), M(n), X(e, n);
	};
	Q(y, (e) => {
		q(r) && e(b);
	});
	var x = V(y, 2), S = (e) => {
		var r = wi(), i = B(r);
		Or(i, 21, () => n().line_items, (e) => e.id, (e, n) => {
			pi(e, {
				get item() {
					return q(n);
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
		}), M(i), M(r), X(e, r);
	};
	Q(x, (e) => {
		n().line_items.length > 0 && e(S);
	}), M(u), H(() => Z(g, `(${n().line_items.length ?? ""})`)), J("click", _, () => t.ondelete?.(n().id)), J("click", v, () => z(r, !q(r))), X(e, u), Pe();
}
dr([
	"keydown",
	"dblclick",
	"click"
]);
//#endregion
//#region src/lib/SubcategoryBlock.svelte
var Di = /* @__PURE__ */ Y("<input type=\"text\" class=\"px-2 py-0.5 border border-white/[0.08] rounded text-sm font-medium text-[var(--color-white)] bg-white/[0.04] focus:ring-2 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] font-[var(--font-ui)]\"/>"), Oi = /* @__PURE__ */ Y("<span role=\"button\" tabindex=\"0\" class=\"font-medium text-[var(--color-concrete)] text-sm font-[var(--font-ui)]\"> </span>"), ki = /* @__PURE__ */ Y("<span class=\"text-xs px-1.5 py-0.5 rounded bg-amber-900/30 text-amber-400 font-medium font-[var(--font-ui)]\">overrides</span>"), Ai = /* @__PURE__ */ Y("<span class=\"text-xs px-1.5 py-0.5 rounded bg-green-900/30 text-green-400 font-medium font-[var(--font-ui)]\"> </span>"), ji = /* @__PURE__ */ Y("<div class=\"text-center\"><span> </span> <div class=\"flex items-center justify-center gap-1 mb-1\"><input type=\"checkbox\" class=\"w-3 h-3 rounded border-white/[0.08] bg-white/[0.04] text-[var(--color-sunburst)] focus:ring-[var(--color-sunburst)]\"/> <span class=\"text-xs text-white/40 font-[var(--font-body)]\"> </span></div> <input type=\"number\" step=\"1\" min=\"0\"/> <div class=\"text-xs text-white/40 mt-0.5 font-mono\"> </div></div>"), Mi = /* @__PURE__ */ Y("<div class=\"bg-[var(--color-granite)] rounded-lg border border-white/[0.08] p-3 mb-3\"><div class=\"flex items-center justify-between mb-2\"><span class=\"text-xs font-semibold text-[var(--color-concrete)] uppercase tracking-wide font-[var(--font-ui)]\">Markup Overrides</span> <button class=\"text-xs text-white/40 hover:text-[var(--color-white)] font-[var(--font-ui)]\">Close</button></div> <div class=\"grid grid-cols-5 gap-2\"></div> <div class=\"mt-3 pt-2 border-t border-white/[0.08] flex items-center gap-2\"><span class=\"text-xs font-medium text-[var(--color-concrete)] font-[var(--font-ui)]\">Lump Sum:</span> <span class=\"text-xs text-white/40\">$</span> <input type=\"number\" step=\"0.01\" min=\"0\" class=\"w-28 text-right text-xs font-mono px-2 py-1 border border-white/[0.08] rounded bg-white/[0.04] text-[var(--color-white)]\n								focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)]\"/> <span class=\"text-xs text-white/40 font-[var(--font-body)]\">added post-markup</span></div></div>"), Ni = /* @__PURE__ */ Y("<span> </span>"), Pi = /* @__PURE__ */ Y("<div class=\"flex items-center gap-3 py-1.5 px-2 bg-amber-900/20 rounded text-xs mb-2 cursor-pointer hover:bg-amber-900/30 transition-colors\" role=\"button\" tabindex=\"0\"><span class=\"font-medium text-amber-400 font-[var(--font-ui)]\">Markup:</span> <!></div>"), Fi = /* @__PURE__ */ Y("<table class=\"w-full\"><thead><tr class=\"text-xs text-white/40 uppercase tracking-wide border-b border-white/[0.06] font-[var(--font-ui)]\"><th class=\"px-1 py-1 text-left w-24\">Type</th><th class=\"px-1 py-1 text-left\">Name</th><th class=\"px-1 py-1 text-right w-20\">Qty</th><th class=\"px-1 py-1 text-center w-16\">Unit</th><th class=\"px-1 py-1 text-right w-24\">Price</th><th class=\"px-2 py-1 text-right w-16\">Markup</th><th class=\"px-2 py-1 text-right w-24\">w/ Markup</th><th class=\"px-2 py-1 text-right w-24\">Total</th><th class=\"w-16\"></th></tr></thead><tbody></tbody></table>"), Ii = /* @__PURE__ */ Y("<div class=\"flex items-center gap-2 ml-4 mt-2\"><input type=\"text\" placeholder=\"Group name\" class=\"flex-1 px-2 py-1 bg-transparent border-0 border-b-2 border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/30\"/> <button class=\"px-2 py-1 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-xs rounded font-[var(--font-ui)] font-semibold hover:brightness-110\">Add</button> <button class=\"px-2 py-1 text-[var(--color-muted-text)] text-xs hover:text-[var(--color-white)] font-[var(--font-ui)]\">Cancel</button></div>"), Li = /* @__PURE__ */ Y("<button class=\"ml-4 mt-2 text-xs text-[var(--color-sunburst)] hover:brightness-110 flex items-center gap-1 font-[var(--font-ui)]\"><svg class=\"w-3 h-3\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Group</button>"), Ri = /* @__PURE__ */ Y("<div class=\"mt-2\"><!></div>"), zi = /* @__PURE__ */ Y("<button class=\"mt-2 text-sm text-[var(--color-sunburst)] hover:brightness-110 flex items-center gap-1 font-[var(--font-ui)]\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Item</button>"), Bi = /* @__PURE__ */ Y("<span class=\"text-green-400\"> </span>"), Vi = /* @__PURE__ */ Y("<div class=\"px-4 pb-3\"><!> <!> <!> <!> <!> <div class=\"flex justify-between items-center mt-2 pt-2 border-t border-white/[0.06] text-sm\"><span class=\"text-[var(--color-muted-text)] font-[var(--font-body)]\"> <!></span> <span class=\"font-mono font-semibold text-[var(--color-white)]\"> </span></div></div>"), Hi = /* @__PURE__ */ Y("<div class=\"border-t border-white/[0.06]\"><div class=\"flex items-center justify-between px-4 py-2 hover:bg-white/[0.03] transition-colors group/subcat\"><button class=\"flex items-center gap-2 text-left flex-1 min-w-0\"><svg fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> <!> <span class=\"text-xs text-white/40 font-[var(--font-body)]\"> </span> <!> <!></button> <div class=\"flex items-center gap-2\"><button class=\"text-xs px-1.5 py-0.5 rounded bg-white/[0.06] text-[var(--color-muted-text)] hover:bg-white/[0.1] transition-colors\" title=\"Configure markup\"><svg class=\"w-3 h-3 inline\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4\"></path></svg></button> <span class=\"font-mono text-sm font-semibold text-[var(--color-white)]\"> </span> <button class=\"opacity-0 group-hover/subcat:opacity-100 text-white/30 hover:text-red-400 transition-opacity p-0.5\" title=\"Delete subcategory\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></button></div></div> <!></div>");
function Ui(e, t) {
	Ne(t, !0);
	let n = ti(t, "subcat", 7), r = ti(t, "collapsed", 3, !1), i = ti(t, "materialsDb", 19, () => []), a = ti(t, "ratesDb", 19, () => []), o = /* @__PURE__ */ R(Bt(r())), s = /* @__PURE__ */ R(!1), c = /* @__PURE__ */ R(!1), l = /* @__PURE__ */ R(Bt(n().name)), u = /* @__PURE__ */ R(!1), d = /* @__PURE__ */ R(""), f = /* @__PURE__ */ L(() => ai(n(), t.globals)), p = /* @__PURE__ */ L(() => [
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
	}))), m = /* @__PURE__ */ R(!1), h = /* @__PURE__ */ L(() => q(p).some((e) => e.isOverride || e.isDisabled)), g = /* @__PURE__ */ L(() => n().line_items.length + n().component_groups.reduce((e, t) => e + t.line_items.length, 0));
	function _(e) {
		t.onsnapshot?.(), n().line_items.push({
			id: bi(),
			category_type: e.category_type,
			item_name: e.item_name,
			quantity: 1,
			unit: e.unit,
			unit_price: e.unit_price,
			is_custom: e.is_custom,
			material_id: e.material_id,
			price_override: !1,
			description: null,
			sort_order: n().line_items.length,
			component_group_id: null
		}), z(s, !1), t.onchange?.();
	}
	function v(e) {
		t.onsnapshot?.();
		let r = n().line_items.findIndex((t) => t.id === e);
		r !== -1 && (n().line_items.splice(r, 1), t.onchange?.());
	}
	function y() {
		z(l, n().name, !0), z(c, !0);
	}
	function b() {
		let e = q(l).trim();
		e && e !== n().name && (t.onsnapshot?.(), n().name = e, t.onchange?.()), z(c, !1);
	}
	function x() {
		t.onsnapshot?.();
		let e = q(d).trim() || "New Group";
		n().component_groups.push({
			id: bi(),
			name: e,
			sort_order: n().component_groups.length,
			line_items: []
		}), z(d, ""), z(u, !1), t.onchange?.();
	}
	function S(e) {
		t.onsnapshot?.();
		let r = n().component_groups.findIndex((t) => t.id === e);
		r !== -1 && (n().component_groups.splice(r, 1), t.onchange?.());
	}
	function ee(e, r) {
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
	function w(e) {
		let r = parseFloat(e.target.value);
		!isNaN(r) && r >= 0 && (n().lump_sum = r, t.onchange?.());
	}
	let te = [
		{
			key: "materials",
			label: "Mat",
			color: "text-blue-400"
		},
		{
			key: "labor",
			label: "Lab",
			color: "text-amber-400"
		},
		{
			key: "equipment",
			label: "Equip",
			color: "text-purple-400"
		},
		{
			key: "subs",
			label: "Subs",
			color: "text-green-400"
		},
		{
			key: "other",
			label: "Other",
			color: "text-[var(--color-concrete)]"
		}
	];
	var T = Hi(), E = B(T), D = B(E), ne = B(D), re = V(ne, 2), O = (e) => {
		var t = Di();
		$(t), nn(t, !0), J("click", t, (e) => e.stopPropagation()), J("keydown", t, (e) => {
			e.key === "Enter" && b(), e.key === "Escape" && z(c, !1);
		}), ur("blur", t, b), Xr(t, () => q(l), (e) => z(l, e)), X(e, t);
	}, ie = (e) => {
		var t = Oi(), r = B(t, !0);
		M(t), H(() => Z(r, n().name)), J("dblclick", t, (e) => {
			e.stopPropagation(), y();
		}), X(e, t);
	};
	Q(re, (e) => {
		q(c) ? e(O) : e(ie, -1);
	});
	var ae = V(re, 2), oe = B(ae);
	M(ae);
	var se = V(ae, 2), ce = (e) => {
		X(e, ki());
	};
	Q(se, (e) => {
		q(h) && e(ce);
	});
	var le = V(se, 2), ue = (e) => {
		var t = Ai(), r = B(t);
		M(t), H((e) => Z(r, `+${e ?? ""} lump sum`), [() => ci(n().lump_sum)]), X(e, t);
	};
	Q(le, (e) => {
		n().lump_sum > 0 && e(ue);
	}), M(D);
	var de = V(D, 2), fe = B(de), pe = V(fe, 2), me = B(pe, !0);
	M(pe);
	var he = V(pe, 2);
	M(de), M(E);
	var ge = V(E, 2), k = (e) => {
		var r = Vi(), o = B(r), c = (e) => {
			var r = Mi(), i = B(r), a = V(B(i), 2);
			M(i);
			var o = V(i, 2);
			Or(o, 21, () => te, wr, (e, r) => {
				let i = /* @__PURE__ */ L(() => q(p).find((e) => e.type === q(r).key));
				var a = ji(), o = B(a), s = B(o, !0);
				M(o);
				var c = V(o, 2), l = B(c);
				$(l);
				var u = V(l, 2), d = B(u, !0);
				M(u), M(c);
				var f = V(c, 2);
				$(f);
				var m = V(f, 2), h = B(m);
				M(m), M(a), H((e) => {
					Ir(o, 1, `block text-xs font-medium ${q(r).color ?? ""} mb-1 font-[var(--font-ui)]`), Z(s, q(r).label), Gr(l, n().markup_enabled[q(r).key]), Z(d, n().markup_enabled[q(r).key] ? "On" : "Off"), Wr(f, n().markup_overrides[q(r).key] ?? ""), Kr(f, "placeholder", `${t.globals[`${q(r).key}_markup`]}%`), f.disabled = !n().markup_enabled[q(r).key], Ir(f, 1, `w-full text-center text-xs font-mono px-1 py-1 border border-white/[0.08] rounded
										focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)]
										${n().markup_enabled[q(r).key] ? "bg-white/[0.04] text-[var(--color-white)]" : "bg-white/[0.02] text-white/20"}
										${q(i)?.isOverride ? "border-amber-500/50 bg-amber-900/20" : ""}`), Z(h, `eff: ${e ?? ""}`);
				}, [() => li(q(i)?.value ?? 0)]), J("change", l, () => C(q(r).key)), J("input", f, (e) => ee(q(r).key, e)), X(e, a);
			}), M(o);
			var s = V(o, 2), c = V(B(s), 4);
			$(c), we(2), M(s), M(r), H(() => Wr(c, n().lump_sum)), J("click", a, () => z(m, !1)), J("input", c, w), X(e, r);
		}, l = (e) => {
			var t = Pi();
			Or(V(B(t), 2), 17, () => q(p), wr, (e, t) => {
				var n = Ni(), r = B(n);
				M(n), H((e) => {
					Ir(n, 1, `${q(t).isDisabled ? "text-white/20 line-through" : q(t).isOverride ? "text-amber-400 font-medium" : "text-white/40"} font-[var(--font-body)]`), Z(r, `${q(t).type ?? ""} ${e ?? ""}`);
				}, [() => li(q(t).value)]), X(e, n);
			}), M(t), J("click", t, () => z(m, !0)), J("keydown", t, (e) => {
				(e.key === "Enter" || e.key === " ") && z(m, !0);
			}), X(e, t);
		};
		Q(o, (e) => {
			q(m) ? e(c) : q(h) && e(l, 1);
		});
		var y = V(o, 2), b = (e) => {
			var r = Fi(), i = V(B(r));
			Or(i, 21, () => n().line_items, (e) => e.id, (e, r) => {
				pi(e, {
					get item() {
						return q(r);
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
					ondelete: v
				});
			}), M(i), M(r), X(e, r);
		};
		Q(y, (e) => {
			(q(g) > 0 || n().line_items.length > 0) && e(b);
		});
		var T = V(y, 2);
		Or(T, 17, () => n().component_groups, (e) => e.id, (e, r) => {
			Ei(e, {
				get group() {
					return q(r);
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
				ondelete: S,
				get materialsDb() {
					return i();
				},
				get ratesDb() {
					return a();
				}
			});
		});
		var E = V(T, 2), D = (e) => {
			var t = Ii(), n = B(t);
			$(n);
			var r = V(n, 2), i = V(r, 2);
			M(t), J("keydown", n, (e) => {
				e.key === "Enter" && x(), e.key === "Escape" && z(u, !1);
			}), Xr(n, () => q(d), (e) => z(d, e)), J("click", r, x), J("click", i, () => z(u, !1)), X(e, t);
		}, ne = (e) => {
			var t = Li();
			J("click", t, () => z(u, !0)), X(e, t);
		};
		Q(E, (e) => {
			q(u) ? e(D) : e(ne, -1);
		});
		var re = V(E, 2), O = (e) => {
			var t = Ri();
			vi(B(t), {
				get materialsDb() {
					return i();
				},
				get ratesDb() {
					return a();
				},
				categoryType: "materials",
				onselect: _,
				oncancel: () => z(s, !1)
			}), M(t), X(e, t);
		}, ie = (e) => {
			var t = zi();
			J("click", t, () => z(s, !0)), X(e, t);
		};
		Q(re, (e) => {
			q(s) ? e(O) : e(ie, -1);
		});
		var ae = V(re, 2), oe = B(ae), se = B(oe), ce = V(se), le = (e) => {
			var t = Bi(), r = B(t);
			M(t), H((e) => Z(r, `+ ${e ?? ""} lump sum`), [() => ci(n().lump_sum)]), X(e, t);
		};
		Q(ce, (e) => {
			n().lump_sum > 0 && e(le);
		}), M(oe);
		var ue = V(oe, 2), de = B(ue, !0);
		M(ue), M(ae), M(r), H((e, t) => {
			Z(se, `Subtotal: ${e ?? ""} `), Z(de, t);
		}, [() => ci(q(f).base), () => ci(q(f).withMarkup)]), X(e, r);
	};
	Q(ge, (e) => {
		q(o) || e(k);
	}), M(T), H((e) => {
		Ir(ne, 0, `w-4 h-4 text-white/40 transition-transform ${q(o) ? "" : "rotate-90"}`), Z(oe, `${q(g) ?? ""} item${q(g) === 1 ? "" : "s"}`), Z(me, e);
	}, [() => ci(q(f).withMarkup)]), J("click", D, () => z(o, !q(o))), J("click", fe, () => z(m, !q(m))), J("click", he, () => t.ondelete?.(n().id)), X(e, T), Pe();
}
dr([
	"click",
	"keydown",
	"dblclick",
	"change",
	"input"
]);
//#endregion
//#region src/lib/SectionBlock.svelte
var Wi = /* @__PURE__ */ Y("<input type=\"text\" class=\"bg-white/[0.06] text-[var(--color-white)] px-2 py-0.5 rounded text-sm font-semibold border border-white/[0.08] focus:ring-2 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] font-[var(--font-ui)]\"/>"), Gi = /* @__PURE__ */ Y("<span role=\"button\" tabindex=\"0\" class=\"font-semibold font-[var(--font-ui)] uppercase tracking-wide\"> </span>"), Ki = /* @__PURE__ */ Y("<div class=\"px-4 py-8 text-center text-[var(--color-muted-text)] text-sm font-[var(--font-body)]\">No subcategories yet</div>"), qi = /* @__PURE__ */ Y("<div class=\"flex items-center gap-2 mt-2\"><input type=\"text\" placeholder=\"Subcategory name\" class=\"flex-1 px-2 py-1.5 bg-transparent border-0 border-b-2 border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/30\"/> <button class=\"px-2 py-1.5 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-xs rounded font-[var(--font-ui)] font-semibold hover:brightness-110\">Add</button> <button class=\"px-2 py-1.5 text-[var(--color-muted-text)] text-xs hover:text-[var(--color-white)] font-[var(--font-ui)]\">Cancel</button></div>"), Ji = /* @__PURE__ */ Y("<button class=\"mt-2 text-sm text-[var(--color-sunburst)] hover:brightness-110 flex items-center gap-1 font-[var(--font-ui)]\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Subcategory</button>"), Yi = /* @__PURE__ */ Y("<!> <div class=\"px-4 pb-3\"><!></div>", 1), Xi = /* @__PURE__ */ Y("<div class=\"mb-4 border border-white/[0.06] rounded-lg overflow-hidden bg-[var(--color-ink)]\"><div class=\"flex items-center justify-between px-4 py-3 bg-[var(--color-granite)] text-[var(--color-white)]\"><button class=\"flex items-center gap-3 hover:bg-white/[0.06] -ml-2 px-2 py-1 rounded transition-colors text-left flex-1 min-w-0\"><svg fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> <!> <span class=\"text-xs text-white/40 font-[var(--font-body)]\"> </span></button> <div class=\"flex items-center gap-2\"><span class=\"font-mono font-semibold\"> </span> <button class=\"text-white/30 hover:text-red-400 transition-colors p-1\" title=\"Delete section\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16\"></path></svg></button></div></div> <!></div>");
function Zi(e, t) {
	Ne(t, !0);
	let n = ti(t, "section", 7), r = ti(t, "collapsed", 3, !1), i = ti(t, "materialsDb", 19, () => []), a = ti(t, "ratesDb", 19, () => []), o = /* @__PURE__ */ R(Bt(r())), s = /* @__PURE__ */ R(!1), c = /* @__PURE__ */ R(Bt(n().name)), l = /* @__PURE__ */ R(!1), u = /* @__PURE__ */ R(""), d = /* @__PURE__ */ L(() => oi(n(), t.globals)), f = /* @__PURE__ */ L(() => n().subcategories.reduce((e, t) => e + t.line_items.length + t.component_groups.reduce((e, t) => e + t.line_items.length, 0), 0));
	function p() {
		z(c, n().name, !0), z(s, !0);
	}
	function m() {
		let e = q(c).trim();
		e && e !== n().name && (t.onsnapshot?.(), n().name = e, t.onchange?.()), z(s, !1);
	}
	function h() {
		t.onsnapshot?.();
		let e = q(u).trim() || "New Subcategory";
		n().subcategories.push({
			id: bi(),
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
		}), z(u, ""), z(l, !1), t.onchange?.();
	}
	function g(e) {
		t.onsnapshot?.();
		let r = n().subcategories.findIndex((t) => t.id === e);
		r !== -1 && (n().subcategories.splice(r, 1), t.onchange?.());
	}
	var _ = Xi(), v = B(_), y = B(v), b = B(y), x = V(b, 2), S = (e) => {
		var t = Wi();
		$(t), nn(t, !0), J("click", t, (e) => e.stopPropagation()), J("keydown", t, (e) => {
			e.key === "Enter" && m(), e.key === "Escape" && z(s, !1);
		}), ur("blur", t, m), Xr(t, () => q(c), (e) => z(c, e)), X(e, t);
	}, ee = (e) => {
		var t = Gi(), r = B(t, !0);
		M(t), H(() => Z(r, n().name)), J("dblclick", t, (e) => {
			e.stopPropagation(), p();
		}), X(e, t);
	};
	Q(x, (e) => {
		q(s) ? e(S) : e(ee, -1);
	});
	var C = V(x, 2), w = B(C);
	M(C), M(y);
	var te = V(y, 2), T = B(te), E = B(T, !0);
	M(T);
	var D = V(T, 2);
	M(te), M(v);
	var ne = V(v, 2), re = (e) => {
		var r = Yi(), o = Zt(r), s = (e) => {
			X(e, Ki());
		}, c = (e) => {
			var r = vr();
			Or(Zt(r), 17, () => n().subcategories, (e) => e.id, (e, n) => {
				Ui(e, {
					get subcat() {
						return q(n);
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
					ondelete: g,
					get materialsDb() {
						return i();
					},
					get ratesDb() {
						return a();
					}
				});
			}), X(e, r);
		};
		Q(o, (e) => {
			n().subcategories.length === 0 && !q(l) ? e(s) : e(c, -1);
		});
		var d = V(o, 2), f = B(d), p = (e) => {
			var t = qi(), n = B(t);
			$(n);
			var r = V(n, 2), i = V(r, 2);
			M(t), J("keydown", n, (e) => {
				e.key === "Enter" && h(), e.key === "Escape" && z(l, !1);
			}), Xr(n, () => q(u), (e) => z(u, e)), J("click", r, h), J("click", i, () => z(l, !1)), X(e, t);
		}, m = (e) => {
			var t = Ji();
			J("click", t, () => z(l, !0)), X(e, t);
		};
		Q(f, (e) => {
			q(l) ? e(p) : e(m, -1);
		}), M(d), X(e, r);
	};
	Q(ne, (e) => {
		q(o) || e(re);
	}), M(_), H((e) => {
		Ir(b, 0, `w-4 h-4 text-white/40 transition-transform ${q(o) ? "" : "rotate-90"}`), Z(w, `${n().subcategories.length ?? ""} subcategor${n().subcategories.length === 1 ? "y" : "ies"}
				· ${q(f) ?? ""} item${q(f) === 1 ? "" : "s"}`), Z(E, e);
	}, [() => ci(q(d).withMarkup)]), J("click", y, () => z(o, !q(o))), J("click", D, () => t.ondelete?.(n().id)), X(e, _), Pe();
}
dr([
	"click",
	"keydown",
	"dblclick"
]);
//#endregion
//#region src/lib/FooterSummary.svelte
var Qi = /* @__PURE__ */ Y("<div class=\"flex flex-col\"><span class=\"text-xs text-white/40 uppercase tracking-wide font-[var(--font-ui)]\"> </span> <span class=\"font-mono\"> </span></div>"), $i = /* @__PURE__ */ Y("<span class=\"text-[var(--color-muted-text)] text-xs font-[var(--font-body)]\">No items yet</span>"), ea = /* @__PURE__ */ Y("<div class=\"fixed bottom-0 left-0 right-0 bg-[var(--color-ink)] text-[var(--color-white)] border-t border-white/[0.08] z-20\"><div class=\"flex items-center justify-between px-4 py-3\"><div class=\"flex items-center gap-6 text-sm\"><div class=\"flex flex-col\"><span class=\"text-xs text-white/40 uppercase tracking-wide font-[var(--font-ui)]\">Base Cost</span> <span class=\"font-mono\"> </span></div> <div class=\"w-px h-8 bg-white/[0.08]\"></div> <!> <!></div> <div class=\"flex flex-col items-end\"><span class=\"text-xs text-white/40 uppercase tracking-wide font-[var(--font-ui)]\">Total</span> <span class=\"font-mono text-lg font-bold text-[var(--color-sunburst)]\"> </span></div></div></div>");
function ta(e, t) {
	Ne(t, !0);
	let n = /* @__PURE__ */ L(() => si(t.estimate)), r = /* @__PURE__ */ L(() => Object.entries(q(n).byType).filter(([, e]) => e > 0).map(([e, t]) => ({
		type: e,
		value: t
	}))), i = {
		materials: "Materials",
		labor: "Labor",
		equipment: "Equipment",
		subs: "Subs",
		other: "Other"
	};
	var a = ea(), o = B(a), s = B(o), c = B(s), l = V(B(c), 2), u = B(l, !0);
	M(l), M(c);
	var d = V(c, 4);
	Or(d, 17, () => q(r), wr, (e, t) => {
		let n = () => q(t).type, r = () => q(t).value;
		var a = Qi(), o = B(a), s = B(o, !0);
		M(o);
		var c = V(o, 2), l = B(c, !0);
		M(c), M(a), H((e) => {
			Z(s, i[n()]), Z(l, e);
		}, [() => ci(r())]), X(e, a);
	});
	var f = V(d, 2), p = (e) => {
		X(e, $i());
	};
	Q(f, (e) => {
		q(r).length === 0 && e(p);
	}), M(s);
	var m = V(s, 2), h = V(B(m), 2), g = B(h, !0);
	M(h), M(m), M(o), M(a), H((e, t) => {
		Z(u, e), Z(g, t);
	}, [() => ci(q(n).base), () => ci(q(n).withMarkup)]), X(e, a), Pe();
}
//#endregion
//#region src/lib/SaveStatus.svelte
var na = /* @__PURE__ */ Y("<button class=\"ml-1 px-1.5 py-0.5 rounded bg-[var(--color-sunburst)]/20 hover:bg-[var(--color-sunburst)]/30 text-[var(--color-sunburst)] transition-colors\" title=\"Save now (Ctrl+S)\">Save</button>"), ra = /* @__PURE__ */ Y("<button class=\"ml-1 px-1.5 py-0.5 rounded bg-red-500/20 hover:bg-red-500/30 text-red-300 transition-colors\" title=\"Retry save\">Retry</button>"), ia = /* @__PURE__ */ Y("<div><span></span> <span> </span> <!> <!></div>");
function aa(e, t) {
	Ne(t, !0);
	let n = /* @__PURE__ */ R(Bt(Date.now()));
	pn(() => {
		let e = setInterval(() => {
			z(n, Date.now(), !0);
		}, 1e4);
		return () => clearInterval(e);
	});
	let r = /* @__PURE__ */ L(() => {
		switch (t.status) {
			case "clean":
				if (t.savedAt) {
					let e = Math.round((q(n) - t.savedAt.getTime()) / 1e3);
					return e < 5 ? "Saved just now" : e < 60 ? `Saved ${e}s ago` : `Saved ${Math.round(e / 60)}m ago`;
				}
				return "Up to date";
			case "dirty": return "Unsaved changes";
			case "saving": return "Saving...";
			case "error": return "Save failed";
			default: return "";
		}
	}), i = /* @__PURE__ */ L(() => {
		switch (t.status) {
			case "clean": return "text-[var(--color-sage)]";
			case "dirty": return "text-[var(--color-sunburst)]";
			case "saving": return "text-blue-400";
			case "error": return "text-red-400";
			default: return "text-white/40";
		}
	}), a = /* @__PURE__ */ L(() => {
		switch (t.status) {
			case "clean": return "bg-[var(--color-sage)]";
			case "dirty": return "bg-[var(--color-sunburst)]";
			case "saving": return "bg-blue-400 animate-pulse";
			case "error": return "bg-red-400 animate-pulse";
			default: return "bg-white/40";
		}
	});
	var o = ia(), s = B(o), c = V(s, 2), l = B(c, !0);
	M(c);
	var u = V(c, 2), d = (e) => {
		var n = na();
		J("click", n, function(...e) {
			t.onsave?.apply(this, e);
		}), X(e, n);
	};
	Q(u, (e) => {
		t.status === "dirty" && t.onsave && e(d);
	});
	var f = V(u, 2), p = (e) => {
		var n = ra();
		J("click", n, function(...e) {
			t.onsave?.apply(this, e);
		}), X(e, n);
	};
	Q(f, (e) => {
		t.status === "error" && t.onsave && e(p);
	}), M(o), H(() => {
		Ir(o, 1, `flex items-center gap-1.5 text-xs font-[var(--font-ui)] ${q(i) ?? ""}`), Ir(s, 1, `w-2 h-2 rounded-full ${q(a) ?? ""}`), Z(l, q(r));
	}), X(e, o), Pe();
}
dr(["click"]);
//#endregion
//#region src/lib/autosave.svelte.js
function oa(e, t = 2e3) {
	let n = /* @__PURE__ */ R("clean"), r = /* @__PURE__ */ R(null), i = null, a = null;
	function o(e) {
		a = e;
	}
	function s() {
		z(n, "dirty"), i && clearTimeout(i), i = setTimeout(() => c(), t);
	}
	async function c() {
		if (i &&= (clearTimeout(i), null), !a) return null;
		let o = a();
		if (!o) return null;
		z(n, "saving");
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
			return z(n, "clean"), z(r, /* @__PURE__ */ new Date(), !0), i;
		} catch (e) {
			return z(n, "error"), console.error("Auto-save failed:", e.message), i = setTimeout(() => c(), t * 2), null;
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
			return q(n);
		},
		get savedAt() {
			return q(r);
		}
	};
}
//#endregion
//#region src/lib/undo.svelte.js
var sa = 20;
function ca() {
	let e = Bt([]), t = /* @__PURE__ */ L(() => e.length > 0);
	function n(t) {
		let n = JSON.parse(JSON.stringify({
			globals: t.globals,
			sections: t.sections
		}));
		e.push(n), e.length > sa && e.shift();
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
			return q(t);
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
var la = /* @__PURE__ */ Y("<div class=\"flex items-center justify-center h-64\"><div class=\"text-[var(--color-muted-text)] font-[var(--font-body)]\">Loading estimate...</div></div>"), ua = /* @__PURE__ */ Y("<div class=\"flex items-center justify-center h-64\"><div class=\"text-red-400 font-[var(--font-body)]\"> </div></div>"), da = /* @__PURE__ */ Y("<label><span class=\"font-medium font-[var(--font-ui)]\"> </span> <input type=\"number\" step=\"1\" min=\"0\"/> <span>%</span></label>"), fa = /* @__PURE__ */ Y("<div class=\"text-center py-16 text-[var(--color-muted-text)]\"><svg class=\"w-12 h-12 mx-auto mb-3 text-white/20\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"1.5\" d=\"M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z\"></path></svg> <p class=\"text-lg font-medium font-[var(--font-ui)]\">No sections yet</p> <p class=\"text-sm mt-1 font-[var(--font-body)]\">Add a section to start building your estimate.</p></div>"), pa = /* @__PURE__ */ Y("<div class=\"flex items-center gap-2 mt-3\"><input type=\"text\" placeholder=\"Section name\" class=\"flex-1 px-3 py-2 bg-transparent border-0 border-b-2 border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/30\"/> <button class=\"px-3 py-2 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-sm rounded font-[var(--font-ui)] font-semibold hover:brightness-110\">Add</button> <button class=\"px-3 py-2 text-[var(--color-muted-text)] text-sm hover:text-[var(--color-white)] font-[var(--font-ui)]\">Cancel</button></div>"), ma = /* @__PURE__ */ Y("<button class=\"mt-3 text-sm text-[var(--color-sunburst)] hover:brightness-110 flex items-center gap-1 font-[var(--font-ui)]\"><svg class=\"w-4 h-4\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Add Section</button>"), ha = /* @__PURE__ */ Y("<div class=\"estimate-builder pb-16\"><div class=\"sticky top-0 z-10 bg-[var(--color-ink)] text-[var(--color-white)] px-4 py-3 flex items-center justify-between border-b border-white/[0.08]\"><div><h1 class=\"text-lg font-semibold font-[var(--font-ui)] uppercase tracking-wide\"> </h1> <span class=\"text-xs text-[var(--color-muted-text)] font-[var(--font-ui)]\">Estimate Builder</span></div> <div class=\"flex items-center gap-3\"><button class=\"text-xs px-2 py-1 rounded bg-white/[0.06] hover:bg-white/[0.1] disabled:opacity-30 disabled:cursor-not-allowed transition-colors flex items-center gap-1 font-[var(--font-ui)]\" title=\"Undo (Ctrl+Z)\"><svg class=\"w-3.5 h-3.5\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M3 10h10a5 5 0 015 5v2M3 10l4-4M3 10l4 4\"></path></svg> Undo</button> <!> <span class=\"text-xs px-2 py-1 rounded bg-white/[0.06] text-[var(--color-concrete)] uppercase tracking-wide font-[var(--font-ui)]\"> </span></div></div> <div class=\"bg-[var(--color-granite)] border-b border-white/[0.06] px-4 py-2\"><div class=\"flex items-center gap-4 text-sm\"><span class=\"font-medium text-[var(--color-concrete)] font-[var(--font-ui)] uppercase tracking-wide text-xs\">Global Markup:</span> <!></div></div> <div class=\"p-4\"><!> <!></div> <!></div>");
function ga(e, t) {
	Ne(t, !0);
	let n = /* @__PURE__ */ R(null), r = /* @__PURE__ */ R(null), i = /* @__PURE__ */ R(!0), a = /* @__PURE__ */ R(!1), o = /* @__PURE__ */ R(""), s = oa(t.projectId), c = ca();
	async function l() {
		try {
			let e = await fetch(`/api/estimate/${t.projectId}`);
			if (!e.ok) throw Error(`Failed to load estimate: ${e.status}`);
			z(n, await e.json(), !0), s.register(() => q(n));
		} catch (e) {
			z(r, e.message, !0);
		} finally {
			z(i, !1);
		}
	}
	pn(() => {
		t.projectId && l();
	}), pn(() => {
		let e = document.getElementById("estimate-root");
		e && (e.dataset.dirty = s.status === "dirty" || s.status === "saving" ? "true" : "false");
	}), pn(() => {
		function e(e) {
			(s.status === "dirty" || s.status === "saving") && e.preventDefault();
		}
		function t(e) {
			if (s.status !== "dirty" && s.status !== "saving") return;
			let t = e.target.closest("a[href]");
			t && t.origin === window.location.origin && (t.closest("#estimate-root") || (e.preventDefault(), confirm("You have unsaved changes. Leave without saving?") && (window.location.href = t.href)));
		}
		function r(e) {
			(e.ctrlKey || e.metaKey) && e.key === "z" && !e.shiftKey && (e.preventDefault(), q(n) && c.undo(q(n)) && u()), (e.ctrlKey || e.metaKey) && e.key === "s" && (e.preventDefault(), s.status === "dirty" && s.save());
		}
		return window.addEventListener("beforeunload", e), document.addEventListener("click", t, !0), window.addEventListener("keydown", r), () => {
			window.removeEventListener("beforeunload", e), document.removeEventListener("click", t, !0), window.removeEventListener("keydown", r), s.destroy();
		};
	});
	function u() {
		s.markDirty();
	}
	function d() {
		q(n) && c.snapshot(q(n));
	}
	function f() {
		d();
		let e = q(o).trim() || "New Section";
		q(n).sections.push({
			id: bi(),
			name: e,
			sort_order: q(n).sections.length,
			subcategories: []
		}), z(o, ""), z(a, !1), u();
	}
	function p(e) {
		d();
		let t = q(n).sections.findIndex((t) => t.id === e);
		t !== -1 && (q(n).sections.splice(t, 1), u());
	}
	function m(e, t) {
		let r = parseFloat(t.target.value);
		!isNaN(r) && r >= 0 && (q(n).globals[`${e}_markup`] = r, u());
	}
	let h = [
		{
			key: "materials",
			label: "Materials",
			bg: "bg-blue-900/30",
			text: "text-blue-400",
			border: "border-blue-800"
		},
		{
			key: "labor",
			label: "Labor",
			bg: "bg-amber-900/30",
			text: "text-amber-400",
			border: "border-amber-800"
		},
		{
			key: "equipment",
			label: "Equipment",
			bg: "bg-purple-900/30",
			text: "text-purple-400",
			border: "border-purple-800"
		},
		{
			key: "subs",
			label: "Subs",
			bg: "bg-green-900/30",
			text: "text-green-400",
			border: "border-green-800"
		},
		{
			key: "other",
			label: "Other",
			bg: "bg-white/[0.06]",
			text: "text-[var(--color-concrete)]",
			border: "border-white/[0.08]"
		}
	];
	var g = vr(), _ = Zt(g), v = (e) => {
		X(e, la());
	}, y = (e) => {
		var t = ua(), n = B(t), i = B(n, !0);
		M(n), M(t), H(() => Z(i, q(r))), X(e, t);
	}, b = (e) => {
		var t = ha(), r = B(t), i = B(r), l = B(i), g = B(l, !0);
		M(l), we(2), M(i);
		var _ = V(i, 2), v = B(_), y = V(v, 2);
		aa(y, {
			get status() {
				return s.status;
			},
			get savedAt() {
				return s.savedAt;
			},
			onsave: () => s.save()
		});
		var b = V(y, 2), x = B(b, !0);
		M(b), M(_), M(r);
		var S = V(r, 2), ee = B(S);
		Or(V(B(ee), 2), 17, () => h, wr, (e, t) => {
			var r = da(), i = B(r), a = B(i, !0);
			M(i);
			var o = V(i, 2);
			$(o), we(2), M(r), H(() => {
				Ir(r, 1, `flex items-center gap-1 px-2 py-0.5 rounded ${q(t).bg ?? ""} ${q(t).text ?? ""} font-mono text-xs`), Z(a, q(t).label), Wr(o, q(n).globals[`${q(t).key}_markup`]), Ir(o, 1, `w-12 text-right bg-transparent border-0 p-0 font-mono text-xs focus:ring-1 focus:ring-[var(--color-sunburst)] rounded ${q(t).text ?? ""}`);
			}), J("input", o, (e) => m(q(t).key, e)), X(e, r);
		}), M(ee), M(S);
		var C = V(S, 2), w = B(C), te = (e) => {
			X(e, fa());
		}, T = (e) => {
			var t = vr();
			Or(Zt(t), 17, () => q(n).sections, (e) => e.id, (e, t) => {
				Zi(e, {
					get section() {
						return q(t);
					},
					get globals() {
						return q(n).globals;
					},
					onchange: u,
					onsnapshot: d,
					ondelete: p,
					get materialsDb() {
						return q(n).materials_db;
					},
					get ratesDb() {
						return q(n).rates_db;
					}
				});
			}), X(e, t);
		};
		Q(w, (e) => {
			q(n).sections.length === 0 && !q(a) ? e(te) : e(T, -1);
		});
		var E = V(w, 2), D = (e) => {
			var t = pa(), n = B(t);
			$(n);
			var r = V(n, 2), i = V(r, 2);
			M(t), J("keydown", n, (e) => {
				e.key === "Enter" && f(), e.key === "Escape" && z(a, !1);
			}), Xr(n, () => q(o), (e) => z(o, e)), J("click", r, f), J("click", i, () => z(a, !1)), X(e, t);
		}, ne = (e) => {
			var t = ma();
			J("click", t, () => z(a, !0)), X(e, t);
		};
		Q(E, (e) => {
			q(a) ? e(D) : e(ne, -1);
		}), M(C), ta(V(C, 2), { get estimate() {
			return q(n);
		} }), M(t), H(() => {
			Z(g, q(n).project.name), v.disabled = !c.canUndo, Z(x, q(n).project.status);
		}), J("click", v, () => {
			c.undo(q(n)) && u();
		}), X(e, t);
	};
	Q(_, (e) => {
		q(i) ? e(v) : q(r) ? e(y, 1) : q(n) && e(b, 2);
	}), X(e, g), Pe();
}
dr([
	"click",
	"input",
	"keydown"
]);
//#endregion
//#region src/main.js
var _a = document.getElementById("estimate-root");
if (_a) {
	let e = _a.dataset.projectId;
	yr(ga, {
		target: _a,
		props: { projectId: e }
	});
}
//#endregion
