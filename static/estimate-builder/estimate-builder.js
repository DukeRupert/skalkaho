//#region node_modules/svelte/src/internal/shared/utils.js
var e = Array.isArray, t = Array.prototype.indexOf, n = Array.prototype.includes, r = Array.from;
Object.keys;
var i = Object.defineProperty, a = Object.getOwnPropertyDescriptor, o = Object.prototype, s = Array.prototype, c = Object.getPrototypeOf, l = Object.isExtensible, u = () => {};
function d(e) {
	for (var t = 0; t < e.length; t++) e[t]();
}
function f() {
	var e, t;
	return {
		promise: new Promise((n, r) => {
			e = n, t = r;
		}),
		resolve: e,
		reject: t
	};
}
var p = 1024, m = 2048, h = 4096, g = 8192, _ = 16384, v = 32768, y = 1 << 25, b = 65536, x = 1 << 19, S = 1 << 20, ee = 1 << 25, te = 65536, C = 1 << 21, ne = 1 << 22, re = 1 << 23, ie = Symbol("$state"), ae = new class extends Error {
	name = "StaleReactionError";
	message = "The reaction that called `getAbortSignal()` was re-run or destroyed";
}();
globalThis.document?.contentType;
//#endregion
//#region node_modules/svelte/src/internal/client/errors.js
function oe() {
	throw Error("https://svelte.dev/e/async_derived_orphan");
}
function se(e, t, n) {
	throw Error("https://svelte.dev/e/each_key_duplicate");
}
function ce(e) {
	throw Error("https://svelte.dev/e/effect_in_teardown");
}
function le() {
	throw Error("https://svelte.dev/e/effect_in_unowned_derived");
}
function ue(e) {
	throw Error("https://svelte.dev/e/effect_orphan");
}
function de() {
	throw Error("https://svelte.dev/e/effect_update_depth_exceeded");
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
var ge = {}, w = Symbol();
function _e(e) {
	console.warn("https://svelte.dev/e/hydration_mismatch");
}
function ve() {
	console.warn("https://svelte.dev/e/svelte_boundary_reset_noop");
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/hydration.js
var T = !1;
function ye(e) {
	T = e;
}
var E;
function D(e) {
	if (e === null) throw _e(), ge;
	return E = e;
}
function be() {
	return D(/* @__PURE__ */ L(E));
}
function O(e) {
	if (T) {
		if (/* @__PURE__ */ L(E) !== null) throw _e(), ge;
		E = e;
	}
}
function xe(e = 1) {
	if (T) {
		for (var t = e, n = E; t--;) n = /* @__PURE__ */ L(n);
		E = n;
	}
}
function Se(e = !0) {
	for (var t = 0, n = E;;) {
		if (n.nodeType === 8) {
			var r = n.data;
			if (r === "]") {
				if (t === 0) return n;
				--t;
			} else (r === "[" || r === "[!" || r[0] === "[" && !isNaN(Number(r.slice(1)))) && (t += 1);
		}
		var i = /* @__PURE__ */ L(n);
		e && n.remove(), n = i;
	}
}
function Ce(e) {
	if (!e || e.nodeType !== 8) throw _e(), ge;
	return e.data;
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/equality.js
function we(e) {
	return e === this.v;
}
function Te(e, t) {
	return e == e ? e !== t || typeof e == "object" && !!e || typeof e == "function" : t == t;
}
function Ee(e) {
	return !Te(e, this.v);
}
//#endregion
//#region node_modules/svelte/src/internal/flags/index.js
var De = !1, Oe = !1, k = null;
function ke(e) {
	k = e;
}
function Ae(e, t = !1, n) {
	k = {
		p: k,
		i: !1,
		c: null,
		e: null,
		s: e,
		x: null,
		r: G,
		l: Oe && !t ? {
			s: null,
			u: null,
			$: []
		} : null
	};
}
function je(e) {
	var t = k, n = t.e;
	if (n !== null) {
		t.e = null;
		for (var r of n) Xt(r);
	}
	return e !== void 0 && (t.x = e), t.i = !0, k = t.p, e ?? {};
}
function Me() {
	return !Oe || k !== null && k.l === null;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/task.js
var Ne = [];
function Pe() {
	var e = Ne;
	Ne = [], d(e);
}
function Fe(e) {
	if (Ne.length === 0 && !Ge) {
		var t = Ne;
		queueMicrotask(() => {
			t === Ne && Pe();
		});
	}
	Ne.push(e);
}
function Ie(e) {
	var t = G;
	if (t === null) return H.f |= re, e;
	if (!(t.f & 32768) && !(t.f & 4)) throw e;
	Le(e, t);
}
function Le(e, t) {
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
var Re = ~(m | h | p);
function A(e, t) {
	e.f = e.f & Re | t;
}
function ze(e) {
	e.f & 512 || e.deps === null ? A(e, p) : A(e, h);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/utils.js
function Be(e) {
	if (e !== null) for (let t of e) !(t.f & 2) || !(t.f & 65536) || (t.f ^= te, Be(t.deps));
}
function Ve(e, t, n) {
	e.f & 2048 ? t.add(e) : e.f & 4096 && n.add(e), Be(e.deps), A(e, p);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/store.js
var He = !1, Ue = /* @__PURE__ */ new Set(), j = null, M = null, We = null, Ge = !1, Ke = !1, qe = null, Je = null, Ye = 0, Xe = 1, Ze = class e {
	id = Xe++;
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
			for (var n of t.d) A(n, m), this.schedule(n);
			for (n of t.m) A(n, h), this.schedule(n);
		}
	}
	#d() {
		Ye++ > 1e3 && Qe();
		let t = this.#a;
		this.#a = [], this.apply();
		var n = qe = [], r = [], i = Je = [];
		for (let e of t) try {
			this.#f(e, n, r);
		} catch (t) {
			throw it(e), t;
		}
		if (j = null, i.length > 0) {
			var a = e.ensure();
			for (let e of i) a.schedule(e);
		}
		if (qe = null, Je = null, this.#u()) {
			this.#p(r), this.#p(n);
			for (let [e, t] of this.#c) rt(e, t);
		} else {
			this.#n === 0 && Ue.delete(this), this.#o.clear(), this.#s.clear();
			for (let e of this.#e) e(this);
			this.#e.clear(), $e(r), $e(n), this.#i?.resolve();
		}
		var o = j;
		if (this.#a.length > 0) {
			let e = o ??= this;
			e.#a.push(...this.#a.filter((t) => !e.#a.includes(t)));
		}
		o !== null && (Ue.add(o), o.#d()), Ue.has(this) || this.#m();
	}
	#f(e, t, n) {
		e.f ^= p;
		for (var r = e.first; r !== null;) {
			var i = r.f, a = (i & 96) != 0;
			if (!(a && i & 1024 || i & 8192 || this.#c.has(r)) && r.fn !== null) {
				a ? r.f ^= p : i & 4 ? t.push(r) : De && i & 16777224 ? n.push(r) : wn(r) && (i & 16 && this.#s.add(r), kn(r));
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
		for (var t = 0; t < e.length; t += 1) Ve(e[t], this.#o, this.#s);
	}
	capture(e, t) {
		t !== w && !this.previous.has(e) && this.previous.set(e, t), e.f & 8388608 || (this.current.set(e, e.v), M?.set(e, e.v));
	}
	activate() {
		j = this;
	}
	deactivate() {
		j = null, M = null;
	}
	flush() {
		try {
			if (Ke = !0, j = this, !this.#u()) {
				for (let e of this.#o) this.#s.delete(e), A(e, m), this.schedule(e);
				for (let e of this.#s) A(e, h), this.schedule(e);
			}
			this.#d();
		} finally {
			Ye = 0, We = null, qe = null, Je = null, Ke = !1, j = null, M = null, Ct.clear();
		}
	}
	discard() {
		for (let e of this.#t) e(this);
		this.#t.clear();
	}
	#m() {
		for (let s of Ue) {
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
					for (var a of t) et(a, n, r, i);
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
		--this.#n, e && --this.#r, !(this.#l || t) && (this.#l = !0, Fe(() => {
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
		return (this.#i ??= f()).promise;
	}
	static ensure() {
		if (j === null) {
			let t = j = new e();
			Ke || (Ue.add(j), Ge || Fe(() => {
				j === t && t.flush();
			}));
		}
		return j;
	}
	apply() {
		if (!De || !this.is_fork && Ue.size === 1) {
			M = null;
			return;
		}
		M = new Map(this.current);
		for (let e of Ue) if (e !== this) for (let [t, n] of e.previous) M.has(t) || M.set(t, n);
	}
	schedule(e) {
		if (We = e, e.b?.is_pending && e.f & 16777228 && !(e.f & 32768)) {
			e.b.defer_effect(e);
			return;
		}
		for (var t = e; t.parent !== null;) {
			t = t.parent;
			var n = t.f;
			if (qe !== null && t === G && (De || (H === null || !(H.f & 2)) && !He)) return;
			if (n & 96) {
				if (!(n & 1024)) return;
				t.f ^= p;
			}
		}
		this.#a.push(t);
	}
};
function Qe() {
	try {
		de();
	} catch (e) {
		Le(e, We);
	}
}
var N = null;
function $e(e) {
	var t = e.length;
	if (t !== 0) {
		for (var n = 0; n < t;) {
			var r = e[n++];
			if (!(r.f & 24576) && wn(r) && (N = /* @__PURE__ */ new Set(), kn(r), r.deps === null && r.first === null && r.nodes === null && r.teardown === null && r.ac === null && sn(r), N?.size > 0)) {
				Ct.clear();
				for (let e of N) {
					if (e.f & 24576) continue;
					let t = [e], n = e.parent;
					for (; n !== null;) N.has(n) && (N.delete(n), t.push(n)), n = n.parent;
					for (let e = t.length - 1; e >= 0; e--) {
						let n = t[e];
						n.f & 24576 || kn(n);
					}
				}
				N.clear();
			}
		}
		N = null;
	}
}
function et(e, t, n, r) {
	if (!n.has(e) && (n.add(e), e.reactions !== null)) for (let i of e.reactions) {
		let e = i.f;
		e & 2 ? et(i, t, n, r) : e & 4194320 && !(e & 2048) && tt(i, t, r) && (A(i, m), nt(i));
	}
}
function tt(e, t, r) {
	let i = r.get(e);
	if (i !== void 0) return i;
	if (e.deps !== null) for (let i of e.deps) {
		if (n.call(t, i)) return !0;
		if (i.f & 2 && tt(i, t, r)) return r.set(i, !0), !0;
	}
	return r.set(e, !1), !1;
}
function nt(e) {
	j.schedule(e);
}
function rt(e, t) {
	if (!(e.f & 32 && e.f & 1024)) {
		e.f & 2048 ? t.d.push(e) : e.f & 4096 && t.m.push(e), A(e, p);
		for (var n = e.first; n !== null;) rt(n, t), n = n.next;
	}
}
function it(e) {
	A(e, p);
	for (var t = e.first; t !== null;) it(t), t = t.next;
}
//#endregion
//#region node_modules/svelte/src/reactivity/create-subscriber.js
function at(e) {
	let t = 0, n = Tt(0), r;
	return () => {
		qt() && (Z(n), $t(() => (t === 0 && (r = Mn(() => e(() => kt(n)))), t += 1, () => {
			Fe(() => {
				--t, t === 0 && (r?.(), r = void 0, kt(n));
			});
		})));
	};
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/boundary.js
var ot = b | x;
function st(e, t, n, r) {
	new ct(e, t, n, r);
}
var ct = class {
	parent;
	is_pending = !1;
	transform_error;
	#e;
	#t = T ? E : null;
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
	#h = at(() => (this.#m = Tt(this.#l), () => {
		this.#m = null;
	}));
	constructor(e, t, n, r) {
		this.#e = e, this.#n = t, this.#r = (e) => {
			var t = G;
			t.b = this, t.f |= 128, n(e);
		}, this.parent = G.b, this.transform_error = r ?? this.parent?.transform_error ?? ((e) => e), this.#i = tn(() => {
			if (T) {
				let e = this.#t;
				be();
				let t = e.data === "[!";
				if (e.data.startsWith("[?")) {
					let t = JSON.parse(e.data.slice(2));
					this.#_(t);
				} else t ? this.#v() : this.#g();
			} else this.#y();
		}, ot), T && (this.#e = E);
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
		e && (this.is_pending = !0, this.#o = B(() => e(this.#e)), Fe(() => {
			var e = this.#c = document.createDocumentFragment(), t = I();
			e.append(t), this.#a = this.#x(() => B(() => this.#r(t))), this.#u === 0 && (this.#e.before(e), this.#c = null, cn(this.#o, () => {
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
				fn(this.#a, e);
				let t = this.#n.pending;
				this.#o = B(() => t(this.#e));
			} else this.#b(j);
		} catch (e) {
			this.error(e);
		}
	}
	#b(e) {
		this.is_pending = !1;
		for (let t of this.#f) A(t, m), e.schedule(t);
		for (let t of this.#p) A(t, h), e.schedule(t);
		this.#f.clear(), this.#p.clear();
	}
	defer_effect(e) {
		Ve(e, this.#f, this.#p);
	}
	is_rendered() {
		return !this.is_pending && (!this.parent || this.parent.is_rendered());
	}
	has_pending_snippet() {
		return !!this.#n.pending;
	}
	#x(e) {
		var t = G, n = H, r = k;
		K(this.#i), W(this.#i), ke(this.#i.ctx);
		try {
			return Ze.ensure(), e();
		} catch (e) {
			return Ie(e), null;
		} finally {
			K(t), W(n), ke(r);
		}
	}
	#S(e, t) {
		if (!this.has_pending_snippet()) {
			this.parent && this.parent.#S(e, t);
			return;
		}
		this.#u += e, this.#u === 0 && (this.#b(t), this.#o && cn(this.#o, () => {
			this.#o = null;
		}), this.#c &&= (this.#e.before(this.#c), null));
	}
	update_pending_count(e, t) {
		this.#S(e, t), this.#l += e, !(!this.#m || this.#d) && (this.#d = !0, Fe(() => {
			this.#d = !1, this.#m && Dt(this.#m, this.#l);
		}));
	}
	get_effect_pending() {
		return this.#h(), Z(this.#m);
	}
	error(e) {
		var t = this.#n.onerror;
		let n = this.#n.failed;
		if (!t && !n) throw e;
		this.#a &&= (V(this.#a), null), this.#o &&= (V(this.#o), null), this.#s &&= (V(this.#s), null), T && (D(this.#t), xe(), D(Se()));
		var r = !1, i = !1;
		let a = () => {
			if (r) {
				ve();
				return;
			}
			r = !0, i && he(), this.#s !== null && cn(this.#s, () => {
				this.#s = null;
			}), this.#x(() => {
				this.#y();
			});
		}, o = (e) => {
			try {
				i = !0, t?.(e, a), i = !1;
			} catch (e) {
				Le(e, this.#i && this.#i.parent);
			}
			n && (this.#s = this.#x(() => {
				try {
					return B(() => {
						var t = G;
						t.b = this, t.f |= 128, n(this.#e, () => e, () => a);
					});
				} catch (e) {
					return Le(e, this.#i.parent), null;
				}
			}));
		};
		Fe(() => {
			var t;
			try {
				t = this.transform_error(e);
			} catch (e) {
				Le(e, this.#i && this.#i.parent);
				return;
			}
			typeof t == "object" && t && typeof t.then == "function" ? t.then(o, (e) => Le(e, this.#i && this.#i.parent)) : o(t);
		});
	}
};
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/async.js
function lt(e, t, n, r) {
	let i = Me() ? pt : ht;
	var a = e.filter((e) => !e.settled);
	if (n.length === 0 && a.length === 0) {
		r(t.map(i));
		return;
	}
	var o = G, s = ut(), c = a.length === 1 ? a[0].promise : a.length > 1 ? Promise.all(a.map((e) => e.promise)) : null;
	function l(e) {
		s();
		try {
			r(e);
		} catch (e) {
			o.f & 16384 || Le(e, o);
		}
		dt();
	}
	if (n.length === 0) {
		c.then(() => l(t.map(i)));
		return;
	}
	var u = ft();
	function d() {
		Promise.all(n.map((e) => /* @__PURE__ */ mt(e))).then((e) => l([...t.map(i), ...e])).catch((e) => Le(e, o)).finally(() => u());
	}
	c ? c.then(() => {
		s(), d(), dt();
	}) : d();
}
function ut() {
	var e = G, t = H, n = k, r = j;
	return function(i = !0) {
		K(e), W(t), ke(n), i && !(e.f & 16384) && (r?.activate(), r?.apply());
	};
}
function dt(e = !0) {
	K(null), W(null), ke(null), e && j?.deactivate();
}
function ft() {
	var e = G.b, t = j, n = e.is_rendered();
	return e.update_pending_count(1, t), t.increment(n), (r = !1) => {
		e.update_pending_count(-1, t), t.decrement(n, r);
	};
}
/* @__NO_SIDE_EFFECTS__ */
function pt(e) {
	var t = 2 | m, n = H !== null && H.f & 2 ? H : null;
	return G !== null && (G.f |= x), {
		ctx: k,
		deps: null,
		effects: null,
		equals: we,
		f: t,
		fn: e,
		reactions: null,
		rv: 0,
		v: w,
		wv: 0,
		parent: n ?? G,
		ac: null
	};
}
/* @__NO_SIDE_EFFECTS__ */
function mt(e, t, n) {
	let r = G;
	r === null && oe();
	var i = void 0, a = Tt(w), o = !H, s = /* @__PURE__ */ new Map();
	return Qt(() => {
		var t = G, n = f();
		i = n.promise;
		try {
			Promise.resolve(e()).then(n.resolve, n.reject).finally(dt);
		} catch (e) {
			n.reject(e), dt();
		}
		var c = j;
		if (o) {
			if (t.f & 32768) var l = ft();
			if (r.b.is_rendered()) s.get(c)?.reject(ae), s.delete(c);
			else {
				for (let e of s.values()) e.reject(ae);
				s.clear();
			}
			s.set(c, n);
		}
		let u = (e, n = void 0) => {
			if (l && l(n === ae), !(n === ae || t.f & 16384)) {
				if (c.activate(), n) a.f |= re, Dt(a, n);
				else {
					a.f & 8388608 && (a.f ^= re), Dt(a, e);
					for (let [e, t] of s) {
						if (s.delete(e), e === c) break;
						t.reject(ae);
					}
				}
				c.deactivate();
			}
		};
		n.promise.then(u, (e) => u(null, e || "unknown"));
	}), Jt(() => {
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
function ht(e) {
	let t = /* @__PURE__ */ pt(e);
	return t.equals = Ee, t;
}
function gt(e) {
	var t = e.effects;
	if (t !== null) {
		e.effects = null;
		for (var n = 0; n < t.length; n += 1) V(t[n]);
	}
}
function _t(e) {
	for (var t = e.parent; t !== null;) {
		if (!(t.f & 2)) return t.f & 16384 ? null : t;
		t = t.parent;
	}
	return null;
}
function vt(e) {
	var t, n = G;
	K(_t(e));
	try {
		e.f &= ~te, gt(e), t = En(e);
	} finally {
		K(n);
	}
	return t;
}
function yt(e) {
	var t = vt(e);
	if (!e.equals(t) && (e.wv = Cn(), (!j?.is_fork || e.deps === null) && (e.v = t, e.deps === null))) {
		A(e, p);
		return;
	}
	hn || (M === null ? ze(e) : (qt() || j?.is_fork) && M.set(e, t));
}
function bt(e) {
	if (e.effects !== null) for (let t of e.effects) (t.teardown || t.ac) && (t.teardown?.(), t.ac?.abort(ae), t.teardown = u, t.ac = null, On(t, 0), rn(t));
}
function xt(e) {
	if (e.effects !== null) for (let t of e.effects) t.teardown && kn(t);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/sources.js
var St = /* @__PURE__ */ new Set(), Ct = /* @__PURE__ */ new Map(), wt = !1;
function Tt(e, t) {
	return {
		f: 0,
		v: e,
		reactions: null,
		equals: we,
		rv: 0,
		wv: 0
	};
}
/* @__NO_SIDE_EFFECTS__ */
function P(e, t) {
	let n = Tt(e, t);
	return _n(n), n;
}
/* @__NO_SIDE_EFFECTS__ */
function Et(e, t = !1, n = !0) {
	let r = Tt(e);
	return t || (r.equals = Ee), Oe && n && k !== null && k.l !== null && (k.l.s ??= []).push(r), r;
}
function F(e, t, r = !1) {
	return H !== null && (!U || H.f & 131072) && Me() && H.f & 4325394 && (q === null || !n.call(q, e)) && me(), Dt(e, r ? jt(t) : t, Je);
}
function Dt(e, t, n = null) {
	if (!e.equals(t)) {
		var r = e.v;
		hn ? Ct.set(e, t) : Ct.set(e, r), e.v = t;
		var i = Ze.ensure();
		if (i.capture(e, r), e.f & 2) {
			let t = e;
			e.f & 2048 && vt(t), ze(t);
		}
		e.wv = Cn(), At(e, m, n), Me() && G !== null && G.f & 1024 && !(G.f & 96) && (X === null ? vn([e]) : X.push(e)), !i.is_fork && St.size > 0 && !wt && Ot();
	}
	return t;
}
function Ot() {
	wt = !1;
	for (let e of St) e.f & 1024 && A(e, h), wn(e) && kn(e);
	St.clear();
}
function kt(e) {
	F(e, e.v + 1);
}
function At(e, t, n) {
	var r = e.reactions;
	if (r !== null) for (var i = Me(), a = r.length, o = 0; o < a; o++) {
		var s = r[o], c = s.f;
		if (!(!i && s === G)) {
			var l = (c & m) === 0;
			if (l && A(s, t), c & 2) {
				var u = s;
				M?.delete(u), c & 65536 || (c & 512 && (s.f |= te), At(u, h, n));
			} else if (l) {
				var d = s;
				c & 16 && N !== null && N.add(d), n === null ? nt(d) : n.push(d);
			}
		}
	}
}
function jt(t) {
	if (typeof t != "object" || !t || ie in t) return t;
	let n = c(t);
	if (n !== o && n !== s) return t;
	var r = /* @__PURE__ */ new Map(), i = e(t), l = /* @__PURE__ */ P(0), u = null, d = xn, f = (e) => {
		if (xn === d) return e();
		var t = H, n = xn;
		W(null), Sn(d);
		var r = e();
		return W(t), Sn(n), r;
	};
	return i && r.set("length", /* @__PURE__ */ P(t.length, u)), new Proxy(t, {
		defineProperty(e, t, n) {
			(!("value" in n) || n.configurable === !1 || n.enumerable === !1 || n.writable === !1) && fe();
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
					let e = f(() => /* @__PURE__ */ P(w, u));
					r.set(t, e), kt(l);
				}
			} else F(n, w), kt(l);
			return !0;
		},
		get(e, n, i) {
			if (n === ie) return t;
			var o = r.get(n), s = n in e;
			if (o === void 0 && (!s || a(e, n)?.writable) && (o = f(() => /* @__PURE__ */ P(jt(s ? e[n] : w), u)), r.set(n, o)), o !== void 0) {
				var c = Z(o);
				return c === w ? void 0 : c;
			}
			return Reflect.get(e, n, i);
		},
		getOwnPropertyDescriptor(e, t) {
			var n = Reflect.getOwnPropertyDescriptor(e, t);
			if (n && "value" in n) {
				var i = r.get(t);
				i && (n.value = Z(i));
			} else if (n === void 0) {
				var a = r.get(t), o = a?.v;
				if (a !== void 0 && o !== w) return {
					enumerable: !0,
					configurable: !0,
					value: o,
					writable: !0
				};
			}
			return n;
		},
		has(e, t) {
			if (t === ie) return !0;
			var n = r.get(t), i = n !== void 0 && n.v !== w || Reflect.has(e, t);
			return (n !== void 0 || G !== null && (!i || a(e, t)?.writable)) && (n === void 0 && (n = f(() => /* @__PURE__ */ P(i ? jt(e[t]) : w, u)), r.set(t, n)), Z(n) === w) ? !1 : i;
		},
		set(e, t, n, o) {
			var s = r.get(t), c = t in e;
			if (i && t === "length") for (var d = n; d < s.v; d += 1) {
				var p = r.get(d + "");
				p === void 0 ? d in e && (p = f(() => /* @__PURE__ */ P(w, u)), r.set(d + "", p)) : F(p, w);
			}
			if (s === void 0) (!c || a(e, t)?.writable) && (s = f(() => /* @__PURE__ */ P(void 0, u)), F(s, jt(n)), r.set(t, s));
			else {
				c = s.v !== w;
				var m = f(() => jt(n));
				F(s, m);
			}
			var h = Reflect.getOwnPropertyDescriptor(e, t);
			if (h?.set && h.set.call(o, n), !c) {
				if (i && typeof t == "string") {
					var g = r.get("length"), _ = Number(t);
					Number.isInteger(_) && _ >= g.v && F(g, _ + 1);
				}
				kt(l);
			}
			return !0;
		},
		ownKeys(e) {
			Z(l);
			var t = Reflect.ownKeys(e).filter((e) => {
				var t = r.get(e);
				return t === void 0 || t.v !== w;
			});
			for (var [n, i] of r) i.v !== w && !(n in e) && t.push(n);
			return t;
		},
		setPrototypeOf() {
			pe();
		}
	});
}
var Mt, Nt, Pt, Ft;
function It() {
	if (Mt === void 0) {
		Mt = window, document, Nt = /Firefox/.test(navigator.userAgent);
		var e = Element.prototype, t = Node.prototype, n = Text.prototype;
		Pt = a(t, "firstChild").get, Ft = a(t, "nextSibling").get, l(e) && (e.__click = void 0, e.__className = void 0, e.__attributes = null, e.__style = void 0, e.__e = void 0), l(n) && (n.__t = void 0);
	}
}
function I(e = "") {
	return document.createTextNode(e);
}
/* @__NO_SIDE_EFFECTS__ */
function Lt(e) {
	return Pt.call(e);
}
/* @__NO_SIDE_EFFECTS__ */
function L(e) {
	return Ft.call(e);
}
function R(e, t) {
	if (!T) return /* @__PURE__ */ Lt(e);
	var n = /* @__PURE__ */ Lt(E);
	if (n === null) n = E.appendChild(I());
	else if (t && n.nodeType !== 3) {
		var r = I();
		return n?.before(r), D(r), r;
	}
	return t && Ht(n), D(n), n;
}
function Rt(e, t = !1) {
	if (!T) {
		var n = /* @__PURE__ */ Lt(e);
		return n instanceof Comment && n.data === "" ? /* @__PURE__ */ L(n) : n;
	}
	if (t) {
		if (E?.nodeType !== 3) {
			var r = I();
			return E?.before(r), D(r), r;
		}
		Ht(E);
	}
	return E;
}
function z(e, t = 1, n = !1) {
	let r = T ? E : e;
	for (var i; t--;) i = r, r = /* @__PURE__ */ L(r);
	if (!T) return r;
	if (n) {
		if (r?.nodeType !== 3) {
			var a = I();
			return r === null ? i?.after(a) : r.before(a), D(a), a;
		}
		Ht(r);
	}
	return D(r), r;
}
function zt(e) {
	e.textContent = "";
}
function Bt() {
	return !De || N !== null ? !1 : (G.f & v) !== 0;
}
function Vt(e, t, n) {
	let r = n ? { is: n } : void 0;
	return document.createElementNS(t ?? "http://www.w3.org/1999/xhtml", e, r);
}
function Ht(e) {
	if (e.nodeValue.length < 65536) return;
	let t = e.nextSibling;
	for (; t !== null && t.nodeType === 3;) t.remove(), e.nodeValue += t.nodeValue, t = e.nextSibling;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/shared.js
function Ut(e) {
	var t = H, n = G;
	W(null), K(null);
	try {
		return e();
	} finally {
		W(t), K(n);
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/effects.js
function Wt(e) {
	G === null && (H === null && ue(e), le()), hn && ce(e);
}
function Gt(e, t) {
	var n = t.last;
	n === null ? t.last = t.first = e : (n.next = e, e.prev = n, t.last = e);
}
function Kt(e, t) {
	var n = G;
	n !== null && n.f & 8192 && (e |= g);
	var r = {
		ctx: k,
		deps: null,
		nodes: null,
		f: e | m | 512,
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
	if (e & 4) qe === null ? Ze.ensure().schedule(r) : qe.push(r);
	else if (t !== null) {
		try {
			kn(r);
		} catch (e) {
			throw V(r), e;
		}
		i.deps === null && i.teardown === null && i.nodes === null && i.first === i.last && !(i.f & 524288) && (i = i.first, e & 16 && e & 65536 && i !== null && (i.f |= b));
	}
	if (i !== null && (i.parent = n, n !== null && Gt(i, n), H !== null && H.f & 2 && !(e & 64))) {
		var a = H;
		(a.effects ??= []).push(i);
	}
	return r;
}
function qt() {
	return H !== null && !U;
}
function Jt(e) {
	let t = Kt(8, null);
	return A(t, p), t.teardown = e, t;
}
function Yt(e) {
	Wt("$effect");
	var t = G.f;
	if (!H && t & 32 && !(t & 32768)) {
		var n = k;
		(n.e ??= []).push(e);
	} else return Xt(e);
}
function Xt(e) {
	return Kt(4 | S, e);
}
function Zt(e) {
	Ze.ensure();
	let t = Kt(64 | x, e);
	return (e = {}) => new Promise((n) => {
		e.outro ? cn(t, () => {
			V(t), n(void 0);
		}) : (V(t), n(void 0));
	});
}
function Qt(e) {
	return Kt(ne | x, e);
}
function $t(e, t = 0) {
	return Kt(8 | t, e);
}
function en(e, t = [], n = [], r = []) {
	lt(r, t, n, (t) => {
		Kt(8, () => e(...t.map(Z)));
	});
}
function tn(e, t = 0) {
	return Kt(16 | t, e);
}
function B(e) {
	return Kt(32 | x, e);
}
function nn(e) {
	var t = e.teardown;
	if (t !== null) {
		let e = hn, n = H;
		gn(!0), W(null);
		try {
			t.call(null);
		} finally {
			gn(e), W(n);
		}
	}
}
function rn(e, t = !1) {
	var n = e.first;
	for (e.first = e.last = null; n !== null;) {
		let e = n.ac;
		e !== null && Ut(() => {
			e.abort(ae);
		});
		var r = n.next;
		n.f & 64 ? n.parent = null : V(n, t), n = r;
	}
}
function an(e) {
	for (var t = e.first; t !== null;) {
		var n = t.next;
		t.f & 32 || V(t), t = n;
	}
}
function V(e, t = !0) {
	var n = !1;
	(t || e.f & 262144) && e.nodes !== null && e.nodes.end !== null && (on(e.nodes.start, e.nodes.end), n = !0), A(e, y), rn(e, t && !n), On(e, 0);
	var r = e.nodes && e.nodes.t;
	if (r !== null) for (let e of r) e.stop();
	nn(e), e.f ^= y, e.f |= _;
	var i = e.parent;
	i !== null && i.first !== null && sn(e), e.next = e.prev = e.teardown = e.ctx = e.deps = e.fn = e.nodes = e.ac = null;
}
function on(e, t) {
	for (; e !== null;) {
		var n = e === t ? null : /* @__PURE__ */ L(e);
		e.remove(), e = n;
	}
}
function sn(e) {
	var t = e.parent, n = e.prev, r = e.next;
	n !== null && (n.next = r), r !== null && (r.prev = n), t !== null && (t.first === e && (t.first = r), t.last === e && (t.last = n));
}
function cn(e, t, n = !0) {
	var r = [];
	ln(e, r, !0);
	var i = () => {
		n && V(e), t && t();
	}, a = r.length;
	if (a > 0) {
		var o = () => --a || i();
		for (var s of r) s.out(o);
	} else i();
}
function ln(e, t, n) {
	if (!(e.f & 8192)) {
		e.f ^= g;
		var r = e.nodes && e.nodes.t;
		if (r !== null) for (let e of r) (e.is_global || n) && t.push(e);
		for (var i = e.first; i !== null;) {
			var a = i.next, o = (i.f & 65536) != 0 || (i.f & 32) != 0 && (e.f & 16) != 0;
			ln(i, t, o ? n : !1), i = a;
		}
	}
}
function un(e) {
	dn(e, !0);
}
function dn(e, t) {
	if (e.f & 8192) {
		e.f ^= g, e.f & 1024 || (A(e, m), Ze.ensure().schedule(e));
		for (var n = e.first; n !== null;) {
			var r = n.next, i = (n.f & 65536) != 0 || (n.f & 32) != 0;
			dn(n, i ? t : !1), n = r;
		}
		var a = e.nodes && e.nodes.t;
		if (a !== null) for (let e of a) (e.is_global || t) && e.in();
	}
}
function fn(e, t) {
	if (e.nodes) for (var n = e.nodes.start, r = e.nodes.end; n !== null;) {
		var i = n === r ? null : /* @__PURE__ */ L(n);
		t.append(n), n = i;
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/legacy.js
var pn = null, mn = !1, hn = !1;
function gn(e) {
	hn = e;
}
var H = null, U = !1;
function W(e) {
	H = e;
}
var G = null;
function K(e) {
	G = e;
}
var q = null;
function _n(e) {
	H !== null && (!De || H.f & 2) && (q === null ? q = [e] : q.push(e));
}
var J = null, Y = 0, X = null;
function vn(e) {
	X = e;
}
var yn = 1, bn = 0, xn = bn;
function Sn(e) {
	xn = e;
}
function Cn() {
	return ++yn;
}
function wn(e) {
	var t = e.f;
	if (t & 2048) return !0;
	if (t & 2 && (e.f &= ~te), t & 4096) {
		for (var n = e.deps, r = n.length, i = 0; i < r; i++) {
			var a = n[i];
			if (wn(a) && yt(a), a.wv > e.wv) return !0;
		}
		t & 512 && M === null && A(e, p);
	}
	return !1;
}
function Tn(e, t, r = !0) {
	var i = e.reactions;
	if (i !== null && !(!De && q !== null && n.call(q, e))) for (var a = 0; a < i.length; a++) {
		var o = i[a];
		o.f & 2 ? Tn(o, t, !1) : t === o && (r ? A(o, m) : o.f & 1024 && A(o, h), nt(o));
	}
}
function En(e) {
	var t = J, n = Y, r = X, i = H, a = q, o = k, s = U, c = xn, l = e.f;
	J = null, Y = 0, X = null, H = l & 96 ? null : e, q = null, ke(e.ctx), U = !1, xn = ++bn, e.ac !== null && (Ut(() => {
		e.ac.abort(ae);
	}), e.ac = null);
	try {
		e.f |= C;
		var u = e.fn, d = u();
		e.f |= v;
		var f = e.deps, p = j?.is_fork;
		if (J !== null) {
			var m;
			if (p || On(e, Y), f !== null && Y > 0) for (f.length = Y + J.length, m = 0; m < J.length; m++) f[Y + m] = J[m];
			else e.deps = f = J;
			if (qt() && e.f & 512) for (m = Y; m < f.length; m++) (f[m].reactions ??= []).push(e);
		} else !p && f !== null && Y < f.length && (On(e, Y), f.length = Y);
		if (Me() && X !== null && !U && f !== null && !(e.f & 6146)) for (m = 0; m < X.length; m++) Tn(X[m], e);
		if (i !== null && i !== e) {
			if (bn++, i.deps !== null) for (let e = 0; e < n; e += 1) i.deps[e].rv = bn;
			if (t !== null) for (let e of t) e.rv = bn;
			X !== null && (r === null ? r = X : r.push(...X));
		}
		return e.f & 8388608 && (e.f ^= re), d;
	} catch (e) {
		return Ie(e);
	} finally {
		e.f ^= C, J = t, Y = n, X = r, H = i, q = a, ke(o), U = s, xn = c;
	}
}
function Dn(e, r) {
	let i = r.reactions;
	if (i !== null) {
		var a = t.call(i, e);
		if (a !== -1) {
			var o = i.length - 1;
			o === 0 ? i = r.reactions = null : (i[a] = i[o], i.pop());
		}
	}
	if (i === null && r.f & 2 && (J === null || !n.call(J, r))) {
		var s = r;
		s.f & 512 && (s.f ^= 512, s.f &= ~te), ze(s), bt(s), On(s, 0);
	}
}
function On(e, t) {
	var n = e.deps;
	if (n !== null) for (var r = t; r < n.length; r++) Dn(e, n[r]);
}
function kn(e) {
	var t = e.f;
	if (!(t & 16384)) {
		A(e, p);
		var n = G, r = mn;
		G = e, mn = !0;
		try {
			t & 16777232 ? an(e) : rn(e), nn(e);
			var i = En(e);
			e.teardown = typeof i == "function" ? i : null, e.wv = yn;
		} finally {
			mn = r, G = n;
		}
	}
}
function Z(e) {
	var t = (e.f & 2) != 0;
	if (pn?.add(e), H !== null && !U && !(G !== null && G.f & 16384) && (q === null || !n.call(q, e))) {
		var r = H.deps;
		if (H.f & 2097152) e.rv < bn && (e.rv = bn, J === null && r !== null && r[Y] === e ? Y++ : J === null ? J = [e] : J.push(e));
		else {
			(H.deps ??= []).push(e);
			var i = e.reactions;
			i === null ? e.reactions = [H] : n.call(i, H) || i.push(H);
		}
	}
	if (hn && Ct.has(e)) return Ct.get(e);
	if (t) {
		var a = e;
		if (hn) {
			var o = a.v;
			return (!(a.f & 1024) && a.reactions !== null || jn(a)) && (o = vt(a)), Ct.set(a, o), o;
		}
		var s = (a.f & 512) == 0 && !U && H !== null && (mn || (H.f & 512) != 0), c = (a.f & v) === 0;
		wn(a) && (s && (a.f |= 512), yt(a)), s && !c && (xt(a), An(a));
	}
	if (M?.has(e)) return M.get(e);
	if (e.f & 8388608) throw e.v;
	return e.v;
}
function An(e) {
	if (e.f |= 512, e.deps !== null) for (let t of e.deps) (t.reactions ??= []).push(e), t.f & 2 && !(t.f & 512) && (xt(t), An(t));
}
function jn(e) {
	if (e.v === w) return !0;
	if (e.deps === null) return !1;
	for (let t of e.deps) if (Ct.has(t) || t.f & 2 && jn(t)) return !0;
	return !1;
}
function Mn(e) {
	var t = U;
	try {
		return U = !0, e();
	} finally {
		U = t;
	}
}
[.../* @__PURE__ */ "allowfullscreen.async.autofocus.autoplay.checked.controls.default.disabled.formnovalidate.indeterminate.inert.ismap.loop.multiple.muted.nomodule.novalidate.open.playsinline.readonly.required.reversed.seamless.selected.webkitdirectory.defer.disablepictureinpicture.disableremoteplayback".split(".")];
var Nn = ["touchstart", "touchmove"];
function Pn(e) {
	return Nn.includes(e);
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/events.js
var Fn = Symbol("events"), In = /* @__PURE__ */ new Set(), Ln = /* @__PURE__ */ new Set(), Rn = null;
function zn(e) {
	var t = this, n = t.ownerDocument, r = e.type, a = e.composedPath?.() || [], o = a[0] || e.target;
	Rn = e;
	var s = 0, c = Rn === e && e[Fn];
	if (c) {
		var l = a.indexOf(c);
		if (l !== -1 && (t === document || t === window)) {
			e[Fn] = t;
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
		var d = H, f = G;
		W(null), K(null);
		try {
			for (var p, m = []; o !== null;) {
				var h = o.assignedSlot || o.parentNode || o.host || null;
				try {
					var g = o[Fn]?.[r];
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
			e[Fn] = t, delete e.currentTarget, W(d), K(f);
		}
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/reconciler.js
var Bn = globalThis?.window?.trustedTypes && /* @__PURE__ */ globalThis.window.trustedTypes.createPolicy("svelte-trusted-html", { createHTML: (e) => e });
function Vn(e) {
	return Bn?.createHTML(e) ?? e;
}
function Hn(e) {
	var t = Vt("template");
	return t.innerHTML = Vn(e.replaceAll("<!>", "<!---->")), t.content;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/template.js
function Un(e, t) {
	var n = G;
	n.nodes === null && (n.nodes = {
		start: e,
		end: t,
		a: null,
		t: null
	});
}
/* @__NO_SIDE_EFFECTS__ */
function Wn(e, t) {
	var n = (t & 1) != 0, r = (t & 2) != 0, i, a = !e.startsWith("<!>");
	return () => {
		if (T) return Un(E, null), E;
		i === void 0 && (i = Hn(a ? e : "<!>" + e), n || (i = /* @__PURE__ */ Lt(i)));
		var t = r || Nt ? document.importNode(i, !0) : i.cloneNode(!0);
		if (n) {
			var o = /* @__PURE__ */ Lt(t), s = t.lastChild;
			Un(o, s);
		} else Un(t, t);
		return t;
	};
}
function Gn() {
	if (T) return Un(E, null), E;
	var e = document.createDocumentFragment(), t = document.createComment(""), n = I();
	return e.append(t, n), Un(t, n), e;
}
function Q(e, t) {
	if (T) {
		var n = G;
		(!(n.f & 32768) || n.nodes.end === null) && (n.nodes.end = E), be();
		return;
	}
	e !== null && e.before(t);
}
function $(e, t) {
	var n = t == null ? "" : typeof t == "object" ? `${t}` : t;
	n !== (e.__t ??= e.nodeValue) && (e.__t = n, e.nodeValue = `${n}`);
}
function Kn(e, t) {
	return Jn(e, t);
}
var qn = /* @__PURE__ */ new Map();
function Jn(e, { target: t, anchor: n, props: i = {}, events: a, context: o, intro: s = !0, transformError: c }) {
	It();
	var l = void 0, u = Zt(() => {
		var s = n ?? t.appendChild(I());
		st(s, { pending: () => {} }, (t) => {
			Ae({});
			var n = k;
			if (o && (n.c = o), a && (i.$$events = a), T && Un(t, null), l = e(t, i) || {}, T && (G.nodes.end = E, E === null || E.nodeType !== 8 || E.data !== "]")) throw _e(), ge;
			je();
		}, c);
		var u = /* @__PURE__ */ new Set(), d = (e) => {
			for (var n = 0; n < e.length; n++) {
				var r = e[n];
				if (!u.has(r)) {
					u.add(r);
					var i = Pn(r);
					for (let e of [t, document]) {
						var a = qn.get(e);
						a === void 0 && (a = /* @__PURE__ */ new Map(), qn.set(e, a));
						var o = a.get(r);
						o === void 0 ? (e.addEventListener(r, zn, { passive: i }), a.set(r, 1)) : a.set(r, o + 1);
					}
				}
			}
		};
		return d(r(In)), Ln.add(d), () => {
			for (var e of u) for (let n of [t, document]) {
				var r = qn.get(n), i = r.get(e);
				--i == 0 ? (n.removeEventListener(e, zn), r.delete(e), r.size === 0 && qn.delete(n)) : r.set(e, i);
			}
			Ln.delete(d), s !== n && s.parentNode?.removeChild(s);
		};
	});
	return Yn.set(l, u), l;
}
var Yn = /* @__PURE__ */ new WeakMap(), Xn = class {
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
			if (n) un(n), this.#r.delete(t);
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
						fn(r, t), t.append(I()), this.#n.set(e, {
							effect: r,
							fragment: t
						});
					} else V(r);
					this.#r.delete(e), this.#t.delete(e);
				};
				this.#i || !n ? (this.#r.add(e), cn(r, i, !1)) : i();
			}
		}
	};
	#o = (e) => {
		this.#e.delete(e);
		let t = Array.from(this.#e.values());
		for (let [e, n] of this.#n) t.includes(e) || (V(n.effect), this.#n.delete(e));
	};
	ensure(e, t) {
		var n = j, r = Bt();
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
		} else T && (this.anchor = E), this.#a(n);
	}
};
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/if.js
function Zn(e, t, n = !1) {
	var r;
	T && (r = E, be());
	var i = new Xn(e), a = n ? b : 0;
	function o(e, t) {
		if (T) {
			var n = Ce(r);
			if (e !== parseInt(n.substring(1))) {
				var a = Se();
				D(a), i.anchor = a, ye(!1), i.ensure(e, t), ye(!0);
				return;
			}
		}
		i.ensure(e, t);
	}
	tn(() => {
		var e = !1;
		t((t, n = 0) => {
			e = !0, o(n, t);
		}), e || o(-1, null);
	}, a);
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/each.js
function Qn(e, t) {
	return t;
}
function $n(e, t, n) {
	for (var i = [], a = t.length, o, s = t.length, c = 0; c < a; c++) {
		let n = t[c];
		cn(n, () => {
			if (o) {
				if (o.pending.delete(n), o.done.add(n), o.pending.size === 0) {
					var t = e.outrogroups;
					er(e, r(o.done)), t.delete(o), t.size === 0 && (e.outrogroups = null);
				}
			} else --s;
		}, !1);
	}
	if (s === 0) {
		var l = i.length === 0 && n !== null;
		if (l) {
			var u = n, d = u.parentNode;
			zt(d), d.append(u), e.items.clear();
		}
		er(e, t, !l);
	} else o = {
		pending: new Set(t),
		done: /* @__PURE__ */ new Set()
	}, (e.outrogroups ??= /* @__PURE__ */ new Set()).add(o);
}
function er(e, t, n = !0) {
	var r;
	if (e.pending.size > 0) {
		r = /* @__PURE__ */ new Set();
		for (let t of e.pending.values()) for (let n of t) r.add(e.items.get(n).e);
	}
	for (var i = 0; i < t.length; i++) {
		var a = t[i];
		r?.has(a) ? (a.f |= ee, fn(a, document.createDocumentFragment())) : V(t[i], n);
	}
}
var tr;
function nr(t, n, i, a, o, s = null) {
	var c = t, l = /* @__PURE__ */ new Map();
	if (n & 4) {
		var u = t;
		c = T ? D(/* @__PURE__ */ Lt(u)) : u.appendChild(I());
	}
	T && be();
	var d = null, f = /* @__PURE__ */ ht(() => {
		var t = i();
		return e(t) ? t : t == null ? [] : r(t);
	}), p, m = /* @__PURE__ */ new Map(), h = !0;
	function g(e) {
		v.effect.f & 16384 || (v.pending.delete(e), v.fallback = d, ir(v, p, c, n, a), d !== null && (p.length === 0 ? d.f & 33554432 ? (d.f ^= ee, or(d, null, c)) : un(d) : cn(d, () => {
			d = null;
		})));
	}
	function _(e) {
		v.pending.delete(e);
	}
	var v = {
		effect: tn(() => {
			p = Z(f);
			var e = p.length;
			let t = !1;
			T && Ce(c) === "[!" != (e === 0) && (c = Se(), D(c), ye(!1), t = !0);
			for (var r = /* @__PURE__ */ new Set(), u = j, v = Bt(), y = 0; y < e; y += 1) {
				T && E.nodeType === 8 && E.data === "]" && (c = E, t = !0, ye(!1));
				var b = p[y], x = a(b, y), S = h ? null : l.get(x);
				S ? (S.v && Dt(S.v, b), S.i && Dt(S.i, y), v && u.unskip_effect(S.e)) : (S = ar(l, h ? c : tr ??= I(), b, x, y, o, n, i), h || (S.e.f |= ee), l.set(x, S)), r.add(x);
			}
			if (e === 0 && s && !d && (h ? d = B(() => s(c)) : (d = B(() => s(tr ??= I())), d.f |= ee)), e > r.size && se("", "", ""), T && e > 0 && D(Se()), !h) if (m.set(u, r), v) {
				for (let [e, t] of l) r.has(e) || u.skip_effect(t.e);
				u.oncommit(g), u.ondiscard(_);
			} else g(u);
			t && ye(!0), Z(f);
		}),
		flags: n,
		items: l,
		pending: m,
		outrogroups: null,
		fallback: d
	};
	h = !1, T && (c = E);
}
function rr(e) {
	for (; e !== null && !(e.f & 32);) e = e.next;
	return e;
}
function ir(e, t, n, i, a) {
	var o = (i & 8) != 0, s = t.length, c = e.items, l = rr(e.effect.first), u, d = null, f, p = [], m = [], h, g, _, v;
	if (o) for (v = 0; v < s; v += 1) h = t[v], g = a(h, v), _ = c.get(g).e, _.f & 33554432 || (_.nodes?.a?.measure(), (f ??= /* @__PURE__ */ new Set()).add(_));
	for (v = 0; v < s; v += 1) {
		if (h = t[v], g = a(h, v), _ = c.get(g).e, e.outrogroups !== null) for (let t of e.outrogroups) t.pending.delete(_), t.done.delete(_);
		if (_.f & 33554432) if (_.f ^= ee, _ === l) or(_, null, n);
		else {
			var y = d ? d.next : l;
			_ === e.effect.last && (e.effect.last = _.prev), _.prev && (_.prev.next = _.next), _.next && (_.next.prev = _.prev), sr(e, d, _), sr(e, _, y), or(_, y, n), d = _, p = [], m = [], l = rr(d.next);
			continue;
		}
		if (_.f & 8192 && (un(_), o && (_.nodes?.a?.unfix(), (f ??= /* @__PURE__ */ new Set()).delete(_))), _ !== l) {
			if (u !== void 0 && u.has(_)) {
				if (p.length < m.length) {
					var b = m[0], x;
					d = b.prev;
					var S = p[0], te = p[p.length - 1];
					for (x = 0; x < p.length; x += 1) or(p[x], b, n);
					for (x = 0; x < m.length; x += 1) u.delete(m[x]);
					sr(e, S.prev, te.next), sr(e, d, S), sr(e, te, b), l = b, d = te, --v, p = [], m = [];
				} else u.delete(_), or(_, l, n), sr(e, _.prev, _.next), sr(e, _, d === null ? e.effect.first : d.next), sr(e, d, _), d = _;
				continue;
			}
			for (p = [], m = []; l !== null && l !== _;) (u ??= /* @__PURE__ */ new Set()).add(l), m.push(l), l = rr(l.next);
			if (l === null) continue;
		}
		_.f & 33554432 || p.push(_), d = _, l = rr(_.next);
	}
	if (e.outrogroups !== null) {
		for (let t of e.outrogroups) t.pending.size === 0 && (er(e, r(t.done)), e.outrogroups?.delete(t));
		e.outrogroups.size === 0 && (e.outrogroups = null);
	}
	if (l !== null || u !== void 0) {
		var C = [];
		if (u !== void 0) for (_ of u) _.f & 8192 || C.push(_);
		for (; l !== null;) !(l.f & 8192) && l !== e.fallback && C.push(l), l = rr(l.next);
		var ne = C.length;
		if (ne > 0) {
			var re = i & 4 && s === 0 ? n : null;
			if (o) {
				for (v = 0; v < ne; v += 1) C[v].nodes?.a?.measure();
				for (v = 0; v < ne; v += 1) C[v].nodes?.a?.fix();
			}
			$n(e, C, re);
		}
	}
	o && Fe(() => {
		if (f !== void 0) for (_ of f) _.nodes?.a?.apply();
	});
}
function ar(e, t, n, r, i, a, o, s) {
	var c = o & 1 ? o & 16 ? Tt(n) : /* @__PURE__ */ Et(n, !1, !1) : null, l = o & 2 ? Tt(i) : null;
	return {
		v: c,
		i: l,
		e: B(() => (a(t, c ?? n, l ?? i, s), () => {
			e.delete(r);
		}))
	};
}
function or(e, t, n) {
	if (e.nodes) for (var r = e.nodes.start, i = e.nodes.end, a = t && !(t.f & 33554432) ? t.nodes.start : n; r !== null;) {
		var o = /* @__PURE__ */ L(r);
		if (a.before(r), r === i) return;
		r = o;
	}
}
function sr(e, t, n) {
	t === null ? e.effect.first = n : t.next = n, n === null ? e.effect.last = t : n.prev = t;
}
//#endregion
//#region node_modules/svelte/src/internal/disclose-version.js
typeof window < "u" && ((window.__svelte ??= {}).v ??= /* @__PURE__ */ new Set()).add("5");
//#endregion
//#region src/EstimateBuilder.svelte
var cr = /* @__PURE__ */ Wn("<div class=\"flex items-center justify-center h-64\"><div class=\"text-slate-500\">Loading estimate...</div></div>"), lr = /* @__PURE__ */ Wn("<div class=\"flex items-center justify-center h-64\"><div class=\"text-red-500\"> </div></div>"), ur = /* @__PURE__ */ Wn("<div class=\"text-center py-16 text-slate-400\"><p class=\"text-lg\">No sections yet</p> <p class=\"text-sm mt-1\">Sections, subcategories, and line items will appear here.</p></div>"), dr = /* @__PURE__ */ Wn("<div class=\"mt-1 text-xs text-slate-400\"> </div>"), fr = /* @__PURE__ */ Wn("<div class=\"ml-4 mt-1 text-xs text-slate-500\"> </div>"), pr = /* @__PURE__ */ Wn("<div class=\"border-t border-slate-200 px-4 py-2\"><div class=\"font-medium text-slate-600 text-sm\"> </div> <!> <!></div>"), mr = /* @__PURE__ */ Wn("<div class=\"mb-4 border border-slate-200 rounded-lg overflow-hidden\"><div class=\"bg-slate-100 px-4 py-2 font-semibold text-slate-700\"> </div> <!></div>"), hr = /* @__PURE__ */ Wn("<div class=\"estimate-builder\"><div class=\"sticky top-0 z-10 bg-slate-900 text-white px-4 py-3 flex items-center justify-between border-b border-slate-700\"><div><h1 class=\"text-lg font-semibold\"> </h1> <span class=\"text-xs text-slate-400\">Estimate Builder</span></div> <div class=\"flex items-center gap-4\"><span class=\"text-xs px-2 py-1 rounded bg-slate-700 text-slate-300 uppercase tracking-wide\"> </span></div></div> <div class=\"bg-slate-50 border-b border-slate-200 px-4 py-2\"><div class=\"flex items-center gap-4 text-sm\"><span class=\"font-medium text-slate-600\">Global Markup:</span> <span class=\"text-slate-500\"> </span> <span class=\"text-slate-500\"> </span> <span class=\"text-slate-500\"> </span> <span class=\"text-slate-500\"> </span> <span class=\"text-slate-500\"> </span></div></div> <div class=\"p-4\"><!></div> <div class=\"px-4 py-2 bg-slate-50 border-t border-slate-200 text-xs text-slate-400\"> </div></div>");
function gr(e, t) {
	Ae(t, !0);
	let n = /* @__PURE__ */ P(null), r = /* @__PURE__ */ P(null), i = /* @__PURE__ */ P(!0);
	async function a() {
		try {
			let e = await fetch(`/api/estimate/${t.projectId}`);
			if (!e.ok) throw Error(`Failed to load estimate: ${e.status}`);
			F(n, await e.json(), !0);
		} catch (e) {
			F(r, e.message, !0);
		} finally {
			F(i, !1);
		}
	}
	Yt(() => {
		t.projectId && a();
	});
	var o = Gn(), s = Rt(o), c = (e) => {
		Q(e, cr());
	}, l = (e) => {
		var t = lr(), n = R(t), i = R(n, !0);
		O(n), O(t), en(() => $(i, Z(r))), Q(e, t);
	}, u = (e) => {
		var t = hr(), r = R(t), i = R(r), a = R(i), o = R(a, !0);
		O(a), xe(2), O(i);
		var s = z(i, 2), c = R(s), l = R(c, !0);
		O(c), O(s), O(r);
		var u = z(r, 2), d = R(u), f = z(R(d), 2), p = R(f);
		O(f);
		var m = z(f, 2), h = R(m);
		O(m);
		var g = z(m, 2), _ = R(g);
		O(g);
		var v = z(g, 2), y = R(v);
		O(v);
		var b = z(v, 2), x = R(b);
		O(b), O(d), O(u);
		var S = z(u, 2), ee = R(S), te = (e) => {
			Q(e, ur());
		}, C = (e) => {
			var t = Gn();
			nr(Rt(t), 17, () => Z(n).sections, Qn, (e, t) => {
				var n = mr(), r = R(n), i = R(r, !0);
				O(r), nr(z(r, 2), 17, () => Z(t).subcategories, Qn, (e, t) => {
					var n = pr(), r = R(n), i = R(r, !0);
					O(r);
					var a = z(r, 2), o = (e) => {
						var n = dr(), r = R(n);
						O(n), en(() => $(r, `${Z(t).line_items.length ?? ""} ungrouped item${Z(t).line_items.length === 1 ? "" : "s"}`)), Q(e, n);
					};
					Zn(a, (e) => {
						Z(t).line_items.length > 0 && e(o);
					}), nr(z(a, 2), 17, () => Z(t).component_groups, Qn, (e, t) => {
						var n = fr(), r = R(n);
						O(n), en(() => $(r, `${Z(t).name ?? ""} (${Z(t).line_items.length ?? ""} item${Z(t).line_items.length === 1 ? "" : "s"})`)), Q(e, n);
					}), O(n), en(() => $(i, Z(t).name)), Q(e, n);
				}), O(n), en(() => $(i, Z(t).name)), Q(e, n);
			}), Q(e, t);
		};
		Zn(ee, (e) => {
			Z(n).sections.length === 0 ? e(te) : e(C, -1);
		}), O(S);
		var ne = z(S, 2), re = R(ne);
		O(ne), O(t), en(() => {
			$(o, Z(n).project.name), $(l, Z(n).project.status), $(p, `Materials ${Z(n).globals.materials_markup ?? ""}%`), $(h, `Labor ${Z(n).globals.labor_markup ?? ""}%`), $(_, `Equipment ${Z(n).globals.equipment_markup ?? ""}%`), $(y, `Subs ${Z(n).globals.subs_markup ?? ""}%`), $(x, `Other ${Z(n).globals.other_markup ?? ""}%`), $(re, `${Z(n).sections.length ?? ""} sections |
			${Z(n).materials_db.length ?? ""} materials |
			${Z(n).rates_db.length ?? ""} rates`);
		}), Q(e, t);
	};
	Zn(s, (e) => {
		Z(i) ? e(c) : Z(r) ? e(l, 1) : Z(n) && e(u, 2);
	}), Q(e, o), je();
}
//#endregion
//#region src/main.js
var _r = document.getElementById("estimate-root");
if (_r) {
	let e = _r.dataset.projectId;
	Kn(gr, {
		target: _r,
		props: { projectId: e }
	});
}
//#endregion
