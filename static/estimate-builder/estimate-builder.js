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
var m = 1024, h = 2048, g = 4096, _ = 8192, v = 16384, y = 32768, b = 1 << 25, x = 65536, S = 1 << 19, ee = 1 << 20, te = 1 << 25, C = 65536, ne = 1 << 21, re = 1 << 22, ie = 1 << 23, ae = Symbol("$state"), oe = Symbol("legacy props"), se = Symbol(""), ce = new class extends Error {
	name = "StaleReactionError";
	message = "The reaction that called `getAbortSignal()` was re-run or destroyed";
}(), le = !!globalThis.document?.contentType && /* @__PURE__ */ globalThis.document.contentType.includes("xml");
//#endregion
//#region node_modules/svelte/src/internal/client/errors.js
function ue() {
	throw Error("https://svelte.dev/e/async_derived_orphan");
}
function de(e, t, n) {
	throw Error("https://svelte.dev/e/each_key_duplicate");
}
function fe(e) {
	throw Error("https://svelte.dev/e/effect_in_teardown");
}
function pe() {
	throw Error("https://svelte.dev/e/effect_in_unowned_derived");
}
function me(e) {
	throw Error("https://svelte.dev/e/effect_orphan");
}
function he() {
	throw Error("https://svelte.dev/e/effect_update_depth_exceeded");
}
function ge(e) {
	throw Error("https://svelte.dev/e/props_invalid_value");
}
function _e() {
	throw Error("https://svelte.dev/e/state_descriptors_fixed");
}
function ve() {
	throw Error("https://svelte.dev/e/state_prototype_fixed");
}
function ye() {
	throw Error("https://svelte.dev/e/state_unsafe_mutation");
}
function be() {
	throw Error("https://svelte.dev/e/svelte_boundary_reset_onerror");
}
//#endregion
//#region node_modules/svelte/src/constants.js
var xe = {}, w = Symbol(), Se = "http://www.w3.org/1999/xhtml";
function Ce(e) {
	console.warn("https://svelte.dev/e/hydration_mismatch");
}
function we() {
	console.warn("https://svelte.dev/e/select_multiple_invalid_value");
}
function Te() {
	console.warn("https://svelte.dev/e/svelte_boundary_reset_noop");
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/hydration.js
var T = !1;
function Ee(e) {
	T = e;
}
var E;
function D(e) {
	if (e === null) throw Ce(), xe;
	return E = e;
}
function De() {
	return D(/* @__PURE__ */ Xt(E));
}
function O(e) {
	if (T) {
		if (/* @__PURE__ */ Xt(E) !== null) throw Ce(), xe;
		E = e;
	}
}
function Oe(e = 1) {
	if (T) {
		for (var t = e, n = E; t--;) n = /* @__PURE__ */ Xt(n);
		E = n;
	}
}
function ke(e = !0) {
	for (var t = 0, n = E;;) {
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
function Ae(e) {
	if (!e || e.nodeType !== 8) throw Ce(), xe;
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
		for (var r of n) fn(r);
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
	if (Be.length === 0 && !et) {
		var t = Be;
		queueMicrotask(() => {
			t === Be && Ve();
		});
	}
	Be.push(e);
}
function Ue(e) {
	var t = W;
	if (t === null) return H.f |= ie, e;
	if (!(t.f & 32768) && !(t.f & 4)) throw e;
	We(e, t);
}
function We(e, t) {
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
var Ge = ~(h | g | m);
function A(e, t) {
	e.f = e.f & Ge | t;
}
function Ke(e) {
	e.f & 512 || e.deps === null ? A(e, m) : A(e, g);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/utils.js
function qe(e) {
	if (e !== null) for (let t of e) !(t.f & 2) || !(t.f & 65536) || (t.f ^= C, qe(t.deps));
}
function Je(e, t, n) {
	e.f & 2048 ? t.add(e) : e.f & 4096 && n.add(e), qe(e.deps), A(e, m);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/store.js
var Ye = !1, Xe = !1;
function Ze(e) {
	var t = Xe;
	try {
		return Xe = !1, [e(), Xe];
	} finally {
		Xe = t;
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/batch.js
var Qe = /* @__PURE__ */ new Set(), j = null, M = null, $e = null, et = !1, tt = !1, nt = null, rt = null, it = 0, at = 1, ot = class e {
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
			for (var n of t.d) A(n, h), this.schedule(n);
			for (n of t.m) A(n, g), this.schedule(n);
		}
	}
	#d() {
		it++ > 1e3 && st();
		let t = this.#a;
		this.#a = [], this.apply();
		var n = nt = [], r = [], i = rt = [];
		for (let e of t) try {
			this.#f(e, n, r);
		} catch (t) {
			throw mt(e), t;
		}
		if (j = null, i.length > 0) {
			var a = e.ensure();
			for (let e of i) a.schedule(e);
		}
		if (nt = null, rt = null, this.#u()) {
			this.#p(r), this.#p(n);
			for (let [e, t] of this.#c) pt(e, t);
		} else {
			this.#n === 0 && Qe.delete(this), this.#o.clear(), this.#s.clear();
			for (let e of this.#e) e(this);
			this.#e.clear(), lt(r), lt(n), this.#i?.resolve();
		}
		var o = j;
		if (this.#a.length > 0) {
			let e = o ??= this;
			e.#a.push(...this.#a.filter((t) => !e.#a.includes(t)));
		}
		o !== null && (Qe.add(o), o.#d()), Qe.has(this) || this.#m();
	}
	#f(e, t, n) {
		e.f ^= m;
		for (var r = e.first; r !== null;) {
			var i = r.f, a = (i & 96) != 0;
			if (!(a && i & 1024 || i & 8192 || this.#c.has(r)) && r.fn !== null) {
				a ? r.f ^= m : i & 4 ? t.push(r) : Pe && i & 16777224 ? n.push(r) : Bn(r) && (i & 16 && this.#s.add(r), Gn(r));
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
		for (var t = 0; t < e.length; t += 1) Je(e[t], this.#o, this.#s);
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
			if (tt = !0, j = this, !this.#u()) {
				for (let e of this.#o) this.#s.delete(e), A(e, h), this.schedule(e);
				for (let e of this.#s) A(e, g), this.schedule(e);
			}
			this.#d();
		} finally {
			it = 0, $e = null, nt = null, rt = null, tt = !1, j = null, M = null, Nt.clear();
		}
	}
	discard() {
		for (let e of this.#t) e(this);
		this.#t.clear();
	}
	#m() {
		for (let s of Qe) {
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
					for (var a of t) ut(a, n, r, i);
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
			tt || (Qe.add(j), et || He(() => {
				j === t && t.flush();
			}));
		}
		return j;
	}
	apply() {
		if (!Pe || !this.is_fork && Qe.size === 1) {
			M = null;
			return;
		}
		M = new Map(this.current);
		for (let e of Qe) if (e !== this) for (let [t, n] of e.previous) M.has(t) || M.set(t, n);
	}
	schedule(e) {
		if ($e = e, e.b?.is_pending && e.f & 16777228 && !(e.f & 32768)) {
			e.b.defer_effect(e);
			return;
		}
		for (var t = e; t.parent !== null;) {
			t = t.parent;
			var n = t.f;
			if (nt !== null && t === W && (Pe || (H === null || !(H.f & 2)) && !Ye)) return;
			if (n & 96) {
				if (!(n & 1024)) return;
				t.f ^= m;
			}
		}
		this.#a.push(t);
	}
};
function st() {
	try {
		he();
	} catch (e) {
		We(e, $e);
	}
}
var ct = null;
function lt(e) {
	var t = e.length;
	if (t !== 0) {
		for (var n = 0; n < t;) {
			var r = e[n++];
			if (!(r.f & 24576) && Bn(r) && (ct = /* @__PURE__ */ new Set(), Gn(r), r.deps === null && r.first === null && r.nodes === null && r.teardown === null && r.ac === null && xn(r), ct?.size > 0)) {
				Nt.clear();
				for (let e of ct) {
					if (e.f & 24576) continue;
					let t = [e], n = e.parent;
					for (; n !== null;) ct.has(n) && (ct.delete(n), t.push(n)), n = n.parent;
					for (let e = t.length - 1; e >= 0; e--) {
						let n = t[e];
						n.f & 24576 || Gn(n);
					}
				}
				ct.clear();
			}
		}
		ct = null;
	}
}
function ut(e, t, n, r) {
	if (!n.has(e) && (n.add(e), e.reactions !== null)) for (let i of e.reactions) {
		let e = i.f;
		e & 2 ? ut(i, t, n, r) : e & 4194320 && !(e & 2048) && dt(i, t, r) && (A(i, h), ft(i));
	}
}
function dt(e, t, r) {
	let i = r.get(e);
	if (i !== void 0) return i;
	if (e.deps !== null) for (let i of e.deps) {
		if (n.call(t, i)) return !0;
		if (i.f & 2 && dt(i, t, r)) return r.set(i, !0), !0;
	}
	return r.set(e, !1), !1;
}
function ft(e) {
	j.schedule(e);
}
function pt(e, t) {
	if (!(e.f & 32 && e.f & 1024)) {
		e.f & 2048 ? t.d.push(e) : e.f & 4096 && t.m.push(e), A(e, m);
		for (var n = e.first; n !== null;) pt(n, t), n = n.next;
	}
}
function mt(e) {
	A(e, m);
	for (var t = e.first; t !== null;) mt(t), t = t.next;
}
//#endregion
//#region node_modules/svelte/src/reactivity/create-subscriber.js
function ht(e) {
	let t = 0, n = Ft(0), r;
	return () => {
		ln() && (Y(n), hn(() => (t === 0 && (r = Jn(() => e(() => zt(n)))), t += 1, () => {
			He(() => {
				--t, t === 0 && (r?.(), r = void 0, zt(n));
			});
		})));
	};
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/boundary.js
var gt = x | S;
function _t(e, t, n, r) {
	new vt(e, t, n, r);
}
var vt = class {
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
	#h = ht(() => (this.#m = Ft(this.#l), () => {
		this.#m = null;
	}));
	constructor(e, t, n, r) {
		this.#e = e, this.#n = t, this.#r = (e) => {
			var t = W;
			t.b = this, t.f |= 128, n(e);
		}, this.parent = W.b, this.transform_error = r ?? this.parent?.transform_error ?? ((e) => e), this.#i = gn(() => {
			if (T) {
				let e = this.#t;
				De();
				let t = e.data === "[!";
				if (e.data.startsWith("[?")) {
					let t = JSON.parse(e.data.slice(2));
					this.#_(t);
				} else t ? this.#v() : this.#g();
			} else this.#y();
		}, gt), T && (this.#e = E);
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
			e.append(t), this.#a = this.#x(() => B(() => this.#r(t))), this.#u === 0 && (this.#e.before(e), this.#c = null, Sn(this.#o, () => {
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
				En(this.#a, e);
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
		Je(e, this.#f, this.#p);
	}
	is_rendered() {
		return !this.is_pending && (!this.parent || this.parent.is_rendered());
	}
	has_pending_snippet() {
		return !!this.#n.pending;
	}
	#x(e) {
		var t = W, n = H, r = k;
		Mn(this.#i), U(this.#i), Ie(this.#i.ctx);
		try {
			return ot.ensure(), e();
		} catch (e) {
			return Ue(e), null;
		} finally {
			Mn(t), U(n), Ie(r);
		}
	}
	#S(e, t) {
		if (!this.has_pending_snippet()) {
			this.parent && this.parent.#S(e, t);
			return;
		}
		this.#u += e, this.#u === 0 && (this.#b(t), this.#o && Sn(this.#o, () => {
			this.#o = null;
		}), this.#c &&= (this.#e.before(this.#c), null));
	}
	update_pending_count(e, t) {
		this.#S(e, t), this.#l += e, !(!this.#m || this.#d) && (this.#d = !0, He(() => {
			this.#d = !1, this.#m && Lt(this.#m, this.#l);
		}));
	}
	get_effect_pending() {
		return this.#h(), Y(this.#m);
	}
	error(e) {
		var t = this.#n.onerror;
		let n = this.#n.failed;
		if (!t && !n) throw e;
		this.#a &&= (V(this.#a), null), this.#o &&= (V(this.#o), null), this.#s &&= (V(this.#s), null), T && (D(this.#t), Oe(), D(ke()));
		var r = !1, i = !1;
		let a = () => {
			if (r) {
				Te();
				return;
			}
			r = !0, i && be(), this.#s !== null && Sn(this.#s, () => {
				this.#s = null;
			}), this.#x(() => {
				this.#y();
			});
		}, o = (e) => {
			try {
				i = !0, t?.(e, a), i = !1;
			} catch (e) {
				We(e, this.#i && this.#i.parent);
			}
			n && (this.#s = this.#x(() => {
				try {
					return B(() => {
						var t = W;
						t.b = this, t.f |= 128, n(this.#e, () => e, () => a);
					});
				} catch (e) {
					return We(e, this.#i.parent), null;
				}
			}));
		};
		He(() => {
			var t;
			try {
				t = this.transform_error(e);
			} catch (e) {
				We(e, this.#i && this.#i.parent);
				return;
			}
			typeof t == "object" && t && typeof t.then == "function" ? t.then(o, (e) => We(e, this.#i && this.#i.parent)) : o(t);
		});
	}
};
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/async.js
function yt(e, t, n, r) {
	let i = ze() ? Ct : Tt;
	var a = e.filter((e) => !e.settled);
	if (n.length === 0 && a.length === 0) {
		r(t.map(i));
		return;
	}
	var o = W, s = bt(), c = a.length === 1 ? a[0].promise : a.length > 1 ? Promise.all(a.map((e) => e.promise)) : null;
	function l(e) {
		s();
		try {
			r(e);
		} catch (e) {
			o.f & 16384 || We(e, o);
		}
		xt();
	}
	if (n.length === 0) {
		c.then(() => l(t.map(i)));
		return;
	}
	var u = St();
	function d() {
		Promise.all(n.map((e) => /* @__PURE__ */ wt(e))).then((e) => l([...t.map(i), ...e])).catch((e) => We(e, o)).finally(() => u());
	}
	c ? c.then(() => {
		s(), d(), xt();
	}) : d();
}
function bt() {
	var e = W, t = H, n = k, r = j;
	return function(i = !0) {
		Mn(e), U(t), Ie(n), i && !(e.f & 16384) && (r?.activate(), r?.apply());
	};
}
function xt(e = !0) {
	Mn(null), U(null), Ie(null), e && j?.deactivate();
}
function St() {
	var e = W.b, t = j, n = e.is_rendered();
	return e.update_pending_count(1, t), t.increment(n), (r = !1) => {
		e.update_pending_count(-1, t), t.decrement(n, r);
	};
}
/* @__NO_SIDE_EFFECTS__ */
function Ct(e) {
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
		v: w,
		wv: 0,
		parent: n ?? W,
		ac: null
	};
}
/* @__NO_SIDE_EFFECTS__ */
function wt(e, t, n) {
	let r = W;
	r === null && ue();
	var i = void 0, a = Ft(w), o = !H, s = /* @__PURE__ */ new Map();
	return mn(() => {
		var t = W, n = p();
		i = n.promise;
		try {
			Promise.resolve(e()).then(n.resolve, n.reject).finally(xt);
		} catch (e) {
			n.reject(e), xt();
		}
		var c = j;
		if (o) {
			if (t.f & 32768) var l = St();
			if (r.b.is_rendered()) s.get(c)?.reject(ce), s.delete(c);
			else {
				for (let e of s.values()) e.reject(ce);
				s.clear();
			}
			s.set(c, n);
		}
		let u = (e, n = void 0) => {
			if (l && l(n === ce), !(n === ce || t.f & 16384)) {
				if (c.activate(), n) a.f |= ie, Lt(a, n);
				else {
					a.f & 8388608 && (a.f ^= ie), Lt(a, e);
					for (let [e, t] of s) {
						if (s.delete(e), e === c) break;
						t.reject(ce);
					}
				}
				c.deactivate();
			}
		};
		n.promise.then(u, (e) => u(null, e || "unknown"));
	}), un(() => {
		for (let e of s.values()) e.reject(ce);
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
	let t = /* @__PURE__ */ Ct(e);
	return Pe || Nn(t), t;
}
/* @__NO_SIDE_EFFECTS__ */
function Tt(e) {
	let t = /* @__PURE__ */ Ct(e);
	return t.equals = Ne, t;
}
function Et(e) {
	var t = e.effects;
	if (t !== null) {
		e.effects = null;
		for (var n = 0; n < t.length; n += 1) V(t[n]);
	}
}
function Dt(e) {
	for (var t = e.parent; t !== null;) {
		if (!(t.f & 2)) return t.f & 16384 ? null : t;
		t = t.parent;
	}
	return null;
}
function Ot(e) {
	var t, n = W;
	Mn(Dt(e));
	try {
		e.f &= ~C, Et(e), t = Hn(e);
	} finally {
		Mn(n);
	}
	return t;
}
function kt(e) {
	var t = Ot(e);
	if (!e.equals(t) && (e.wv = zn(), (!j?.is_fork || e.deps === null) && (e.v = t, e.deps === null))) {
		A(e, m);
		return;
	}
	kn || (M === null ? Ke(e) : (ln() || j?.is_fork) && M.set(e, t));
}
function At(e) {
	if (e.effects !== null) for (let t of e.effects) (t.teardown || t.ac) && (t.teardown?.(), t.ac?.abort(ce), t.teardown = d, t.ac = null, Wn(t, 0), vn(t));
}
function jt(e) {
	if (e.effects !== null) for (let t of e.effects) t.teardown && Gn(t);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/sources.js
var Mt = /* @__PURE__ */ new Set(), Nt = /* @__PURE__ */ new Map(), Pt = !1;
function Ft(e, t) {
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
	let n = Ft(e, t);
	return Nn(n), n;
}
/* @__NO_SIDE_EFFECTS__ */
function It(e, t = !1, n = !0) {
	let r = Ft(e);
	return t || (r.equals = Ne), Fe && n && k !== null && k.l !== null && (k.l.s ??= []).push(r), r;
}
function F(e, t, r = !1) {
	return H !== null && (!jn || H.f & 131072) && ze() && H.f & 4325394 && (G === null || !n.call(G, e)) && ye(), Lt(e, r ? Vt(t) : t, rt);
}
function Lt(e, t, n = null) {
	if (!e.equals(t)) {
		var r = e.v;
		kn ? Nt.set(e, t) : Nt.set(e, r), e.v = t;
		var i = ot.ensure();
		if (i.capture(e, r), e.f & 2) {
			let t = e;
			e.f & 2048 && Ot(t), Ke(t);
		}
		e.wv = zn(), Bt(e, h, n), ze() && W !== null && W.f & 1024 && !(W.f & 96) && (J === null ? Pn([e]) : J.push(e)), !i.is_fork && Mt.size > 0 && !Pt && Rt();
	}
	return t;
}
function Rt() {
	Pt = !1;
	for (let e of Mt) e.f & 1024 && A(e, g), Bn(e) && Gn(e);
	Mt.clear();
}
function zt(e) {
	F(e, e.v + 1);
}
function Bt(e, t, n) {
	var r = e.reactions;
	if (r !== null) for (var i = ze(), a = r.length, o = 0; o < a; o++) {
		var s = r[o], c = s.f;
		if (!(!i && s === W)) {
			var l = (c & h) === 0;
			if (l && A(s, t), c & 2) {
				var u = s;
				M?.delete(u), c & 65536 || (c & 512 && (s.f |= C), Bt(u, g, n));
			} else if (l) {
				var d = s;
				c & 16 && ct !== null && ct.add(d), n === null ? ft(d) : n.push(d);
			}
		}
	}
}
function Vt(t) {
	if (typeof t != "object" || !t || ae in t) return t;
	let n = l(t);
	if (n !== s && n !== c) return t;
	var r = /* @__PURE__ */ new Map(), i = e(t), o = /* @__PURE__ */ P(0), u = null, d = Ln, f = (e) => {
		if (Ln === d) return e();
		var t = H, n = Ln;
		U(null), Rn(d);
		var r = e();
		return U(t), Rn(n), r;
	};
	return i && r.set("length", /* @__PURE__ */ P(t.length, u)), new Proxy(t, {
		defineProperty(e, t, n) {
			(!("value" in n) || n.configurable === !1 || n.enumerable === !1 || n.writable === !1) && _e();
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
					r.set(t, e), zt(o);
				}
			} else F(n, w), zt(o);
			return !0;
		},
		get(e, n, i) {
			if (n === ae) return t;
			var o = r.get(n), s = n in e;
			if (o === void 0 && (!s || a(e, n)?.writable) && (o = f(() => /* @__PURE__ */ P(Vt(s ? e[n] : w), u)), r.set(n, o)), o !== void 0) {
				var c = Y(o);
				return c === w ? void 0 : c;
			}
			return Reflect.get(e, n, i);
		},
		getOwnPropertyDescriptor(e, t) {
			var n = Reflect.getOwnPropertyDescriptor(e, t);
			if (n && "value" in n) {
				var i = r.get(t);
				i && (n.value = Y(i));
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
			if (t === ae) return !0;
			var n = r.get(t), i = n !== void 0 && n.v !== w || Reflect.has(e, t);
			return (n !== void 0 || W !== null && (!i || a(e, t)?.writable)) && (n === void 0 && (n = f(() => /* @__PURE__ */ P(i ? Vt(e[t]) : w, u)), r.set(t, n)), Y(n) === w) ? !1 : i;
		},
		set(e, t, n, s) {
			var c = r.get(t), l = t in e;
			if (i && t === "length") for (var d = n; d < c.v; d += 1) {
				var p = r.get(d + "");
				p === void 0 ? d in e && (p = f(() => /* @__PURE__ */ P(w, u)), r.set(d + "", p)) : F(p, w);
			}
			if (c === void 0) (!l || a(e, t)?.writable) && (c = f(() => /* @__PURE__ */ P(void 0, u)), F(c, Vt(n)), r.set(t, c));
			else {
				l = c.v !== w;
				var m = f(() => Vt(n));
				F(c, m);
			}
			var h = Reflect.getOwnPropertyDescriptor(e, t);
			if (h?.set && h.set.call(s, n), !l) {
				if (i && typeof t == "string") {
					var g = r.get("length"), _ = Number(t);
					Number.isInteger(_) && _ >= g.v && F(g, _ + 1);
				}
				zt(o);
			}
			return !0;
		},
		ownKeys(e) {
			Y(o);
			var t = Reflect.ownKeys(e).filter((e) => {
				var t = r.get(e);
				return t === void 0 || t.v !== w;
			});
			for (var [n, i] of r) i.v !== w && !(n in e) && t.push(n);
			return t;
		},
		setPrototypeOf() {
			ve();
		}
	});
}
function Ht(e) {
	try {
		if (typeof e == "object" && e && ae in e) return e[ae];
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
function I(e = "") {
	return document.createTextNode(e);
}
/* @__NO_SIDE_EFFECTS__ */
function Yt(e) {
	return Kt.call(e);
}
/* @__NO_SIDE_EFFECTS__ */
function Xt(e) {
	return qt.call(e);
}
function L(e, t) {
	if (!T) return /* @__PURE__ */ Yt(e);
	var n = /* @__PURE__ */ Yt(E);
	if (n === null) n = E.appendChild(I());
	else if (t && n.nodeType !== 3) {
		var r = I();
		return n?.before(r), D(r), r;
	}
	return t && tn(n), D(n), n;
}
function Zt(e, t = !1) {
	if (!T) {
		var n = /* @__PURE__ */ Yt(e);
		return n instanceof Comment && n.data === "" ? /* @__PURE__ */ Xt(n) : n;
	}
	if (t) {
		if (E?.nodeType !== 3) {
			var r = I();
			return E?.before(r), D(r), r;
		}
		tn(E);
	}
	return E;
}
function R(e, t = 1, n = !1) {
	let r = T ? E : e;
	for (var i; t--;) i = r, r = /* @__PURE__ */ Xt(r);
	if (!T) return r;
	if (n) {
		if (r?.nodeType !== 3) {
			var a = I();
			return r === null ? i?.after(a) : r.before(a), D(a), a;
		}
		tn(r);
	}
	return D(r), r;
}
function Qt(e) {
	e.textContent = "";
}
function $t() {
	return !Pe || ct !== null ? !1 : (W.f & y) !== 0;
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
var nn = !1;
function rn() {
	nn || (nn = !0, document.addEventListener("reset", (e) => {
		Promise.resolve().then(() => {
			if (!e.defaultPrevented) for (let t of e.target.elements) t.__on_r?.();
		});
	}, { capture: !0 }));
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/shared.js
function an(e) {
	var t = H, n = W;
	U(null), Mn(null);
	try {
		return e();
	} finally {
		U(t), Mn(n);
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/effects.js
function on(e) {
	W === null && (H === null && me(e), pe()), kn && fe(e);
}
function sn(e, t) {
	var n = t.last;
	n === null ? t.last = t.first = e : (n.next = e, e.prev = n, t.last = e);
}
function cn(e, t) {
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
	if (e & 4) nt === null ? ot.ensure().schedule(r) : nt.push(r);
	else if (t !== null) {
		try {
			Gn(r);
		} catch (e) {
			throw V(r), e;
		}
		i.deps === null && i.teardown === null && i.nodes === null && i.first === i.last && !(i.f & 524288) && (i = i.first, e & 16 && e & 65536 && i !== null && (i.f |= x));
	}
	if (i !== null && (i.parent = n, n !== null && sn(i, n), H !== null && H.f & 2 && !(e & 64))) {
		var a = H;
		(a.effects ??= []).push(i);
	}
	return r;
}
function ln() {
	return H !== null && !jn;
}
function un(e) {
	let t = cn(8, null);
	return A(t, m), t.teardown = e, t;
}
function dn(e) {
	on("$effect");
	var t = W.f;
	if (!H && t & 32 && !(t & 32768)) {
		var n = k;
		(n.e ??= []).push(e);
	} else return fn(e);
}
function fn(e) {
	return cn(4 | ee, e);
}
function pn(e) {
	ot.ensure();
	let t = cn(64 | S, e);
	return (e = {}) => new Promise((n) => {
		e.outro ? Sn(t, () => {
			V(t), n(void 0);
		}) : (V(t), n(void 0));
	});
}
function mn(e) {
	return cn(re | S, e);
}
function hn(e, t = 0) {
	return cn(8 | t, e);
}
function z(e, t = [], n = [], r = []) {
	yt(r, t, n, (t) => {
		cn(8, () => e(...t.map(Y)));
	});
}
function gn(e, t = 0) {
	return cn(16 | t, e);
}
function B(e) {
	return cn(32 | S, e);
}
function _n(e) {
	var t = e.teardown;
	if (t !== null) {
		let e = kn, n = H;
		An(!0), U(null);
		try {
			t.call(null);
		} finally {
			An(e), U(n);
		}
	}
}
function vn(e, t = !1) {
	var n = e.first;
	for (e.first = e.last = null; n !== null;) {
		let e = n.ac;
		e !== null && an(() => {
			e.abort(ce);
		});
		var r = n.next;
		n.f & 64 ? n.parent = null : V(n, t), n = r;
	}
}
function yn(e) {
	for (var t = e.first; t !== null;) {
		var n = t.next;
		t.f & 32 || V(t), t = n;
	}
}
function V(e, t = !0) {
	var n = !1;
	(t || e.f & 262144) && e.nodes !== null && e.nodes.end !== null && (bn(e.nodes.start, e.nodes.end), n = !0), A(e, b), vn(e, t && !n), Wn(e, 0);
	var r = e.nodes && e.nodes.t;
	if (r !== null) for (let e of r) e.stop();
	_n(e), e.f ^= b, e.f |= v;
	var i = e.parent;
	i !== null && i.first !== null && xn(e), e.next = e.prev = e.teardown = e.ctx = e.deps = e.fn = e.nodes = e.ac = null;
}
function bn(e, t) {
	for (; e !== null;) {
		var n = e === t ? null : /* @__PURE__ */ Xt(e);
		e.remove(), e = n;
	}
}
function xn(e) {
	var t = e.parent, n = e.prev, r = e.next;
	n !== null && (n.next = r), r !== null && (r.prev = n), t !== null && (t.first === e && (t.first = r), t.last === e && (t.last = n));
}
function Sn(e, t, n = !0) {
	var r = [];
	Cn(e, r, !0);
	var i = () => {
		n && V(e), t && t();
	}, a = r.length;
	if (a > 0) {
		var o = () => --a || i();
		for (var s of r) s.out(o);
	} else i();
}
function Cn(e, t, n) {
	if (!(e.f & 8192)) {
		e.f ^= _;
		var r = e.nodes && e.nodes.t;
		if (r !== null) for (let e of r) (e.is_global || n) && t.push(e);
		for (var i = e.first; i !== null;) {
			var a = i.next, o = (i.f & 65536) != 0 || (i.f & 32) != 0 && (e.f & 16) != 0;
			Cn(i, t, o ? n : !1), i = a;
		}
	}
}
function wn(e) {
	Tn(e, !0);
}
function Tn(e, t) {
	if (e.f & 8192) {
		e.f ^= _, e.f & 1024 || (A(e, h), ot.ensure().schedule(e));
		for (var n = e.first; n !== null;) {
			var r = n.next, i = (n.f & 65536) != 0 || (n.f & 32) != 0;
			Tn(n, i ? t : !1), n = r;
		}
		var a = e.nodes && e.nodes.t;
		if (a !== null) for (let e of a) (e.is_global || t) && e.in();
	}
}
function En(e, t) {
	if (e.nodes) for (var n = e.nodes.start, r = e.nodes.end; n !== null;) {
		var i = n === r ? null : /* @__PURE__ */ Xt(n);
		t.append(n), n = i;
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/legacy.js
var Dn = null, On = !1, kn = !1;
function An(e) {
	kn = e;
}
var H = null, jn = !1;
function U(e) {
	H = e;
}
var W = null;
function Mn(e) {
	W = e;
}
var G = null;
function Nn(e) {
	H !== null && (!Pe || H.f & 2) && (G === null ? G = [e] : G.push(e));
}
var K = null, q = 0, J = null;
function Pn(e) {
	J = e;
}
var Fn = 1, In = 0, Ln = In;
function Rn(e) {
	Ln = e;
}
function zn() {
	return ++Fn;
}
function Bn(e) {
	var t = e.f;
	if (t & 2048) return !0;
	if (t & 2 && (e.f &= ~C), t & 4096) {
		for (var n = e.deps, r = n.length, i = 0; i < r; i++) {
			var a = n[i];
			if (Bn(a) && kt(a), a.wv > e.wv) return !0;
		}
		t & 512 && M === null && A(e, m);
	}
	return !1;
}
function Vn(e, t, r = !0) {
	var i = e.reactions;
	if (i !== null && !(!Pe && G !== null && n.call(G, e))) for (var a = 0; a < i.length; a++) {
		var o = i[a];
		o.f & 2 ? Vn(o, t, !1) : t === o && (r ? A(o, h) : o.f & 1024 && A(o, g), ft(o));
	}
}
function Hn(e) {
	var t = K, n = q, r = J, i = H, a = G, o = k, s = jn, c = Ln, l = e.f;
	K = null, q = 0, J = null, H = l & 96 ? null : e, G = null, Ie(e.ctx), jn = !1, Ln = ++In, e.ac !== null && (an(() => {
		e.ac.abort(ce);
	}), e.ac = null);
	try {
		e.f |= ne;
		var u = e.fn, d = u();
		e.f |= y;
		var f = e.deps, p = j?.is_fork;
		if (K !== null) {
			var m;
			if (p || Wn(e, q), f !== null && q > 0) for (f.length = q + K.length, m = 0; m < K.length; m++) f[q + m] = K[m];
			else e.deps = f = K;
			if (ln() && e.f & 512) for (m = q; m < f.length; m++) (f[m].reactions ??= []).push(e);
		} else !p && f !== null && q < f.length && (Wn(e, q), f.length = q);
		if (ze() && J !== null && !jn && f !== null && !(e.f & 6146)) for (m = 0; m < J.length; m++) Vn(J[m], e);
		if (i !== null && i !== e) {
			if (In++, i.deps !== null) for (let e = 0; e < n; e += 1) i.deps[e].rv = In;
			if (t !== null) for (let e of t) e.rv = In;
			J !== null && (r === null ? r = J : r.push(...J));
		}
		return e.f & 8388608 && (e.f ^= ie), d;
	} catch (e) {
		return Ue(e);
	} finally {
		e.f ^= ne, K = t, q = n, J = r, H = i, G = a, Ie(o), jn = s, Ln = c;
	}
}
function Un(e, r) {
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
		s.f & 512 && (s.f ^= 512, s.f &= ~C), Ke(s), At(s), Wn(s, 0);
	}
}
function Wn(e, t) {
	var n = e.deps;
	if (n !== null) for (var r = t; r < n.length; r++) Un(e, n[r]);
}
function Gn(e) {
	var t = e.f;
	if (!(t & 16384)) {
		A(e, m);
		var n = W, r = On;
		W = e, On = !0;
		try {
			t & 16777232 ? yn(e) : vn(e), _n(e);
			var i = Hn(e);
			e.teardown = typeof i == "function" ? i : null, e.wv = Fn;
		} finally {
			On = r, W = n;
		}
	}
}
function Y(e) {
	var t = (e.f & 2) != 0;
	if (Dn?.add(e), H !== null && !jn && !(W !== null && W.f & 16384) && (G === null || !n.call(G, e))) {
		var r = H.deps;
		if (H.f & 2097152) e.rv < In && (e.rv = In, K === null && r !== null && r[q] === e ? q++ : K === null ? K = [e] : K.push(e));
		else {
			(H.deps ??= []).push(e);
			var i = e.reactions;
			i === null ? e.reactions = [H] : n.call(i, H) || i.push(H);
		}
	}
	if (kn && Nt.has(e)) return Nt.get(e);
	if (t) {
		var a = e;
		if (kn) {
			var o = a.v;
			return (!(a.f & 1024) && a.reactions !== null || qn(a)) && (o = Ot(a)), Nt.set(a, o), o;
		}
		var s = (a.f & 512) == 0 && !jn && H !== null && (On || (H.f & 512) != 0), c = (a.f & y) === 0;
		Bn(a) && (s && (a.f |= 512), kt(a)), s && !c && (jt(a), Kn(a));
	}
	if (M?.has(e)) return M.get(e);
	if (e.f & 8388608) throw e.v;
	return e.v;
}
function Kn(e) {
	if (e.f |= 512, e.deps !== null) for (let t of e.deps) (t.reactions ??= []).push(e), t.f & 2 && !(t.f & 512) && (jt(t), Kn(t));
}
function qn(e) {
	if (e.v === w) return !0;
	if (e.deps === null) return !1;
	for (let t of e.deps) if (Nt.has(t) || t.f & 2 && qn(t)) return !0;
	return !1;
}
function Jn(e) {
	var t = jn;
	try {
		return jn = !0, e();
	} finally {
		jn = t;
	}
}
[.../* @__PURE__ */ "allowfullscreen.async.autofocus.autoplay.checked.controls.default.disabled.formnovalidate.indeterminate.inert.ismap.loop.multiple.muted.nomodule.novalidate.open.playsinline.readonly.required.reversed.seamless.selected.webkitdirectory.defer.disablepictureinpicture.disableremoteplayback".split(".")];
var Yn = ["touchstart", "touchmove"];
function Xn(e) {
	return Yn.includes(e);
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/events.js
var Zn = Symbol("events"), Qn = /* @__PURE__ */ new Set(), $n = /* @__PURE__ */ new Set();
function er(e, t, n) {
	(t[Zn] ??= {})[e] = n;
}
function tr(e) {
	for (var t = 0; t < e.length; t++) Qn.add(e[t]);
	for (var n of $n) n(e);
}
var nr = null;
function rr(e) {
	var t = this, n = t.ownerDocument, r = e.type, a = e.composedPath?.() || [], o = a[0] || e.target;
	nr = e;
	var s = 0, c = nr === e && e[Zn];
	if (c) {
		var l = a.indexOf(c);
		if (l !== -1 && (t === document || t === window)) {
			e[Zn] = t;
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
		U(null), Mn(null);
		try {
			for (var p, m = []; o !== null;) {
				var h = o.assignedSlot || o.parentNode || o.host || null;
				try {
					var g = o[Zn]?.[r];
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
			e[Zn] = t, delete e.currentTarget, U(d), Mn(f);
		}
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/reconciler.js
var ir = globalThis?.window?.trustedTypes && /* @__PURE__ */ globalThis.window.trustedTypes.createPolicy("svelte-trusted-html", { createHTML: (e) => e });
function ar(e) {
	return ir?.createHTML(e) ?? e;
}
function or(e) {
	var t = en("template");
	return t.innerHTML = ar(e.replaceAll("<!>", "<!---->")), t.content;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/template.js
function sr(e, t) {
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
		if (T) return sr(E, null), E;
		i === void 0 && (i = or(a ? e : "<!>" + e), n || (i = /* @__PURE__ */ Yt(i)));
		var t = r || Gt ? document.importNode(i, !0) : i.cloneNode(!0);
		if (n) {
			var o = /* @__PURE__ */ Yt(t), s = t.lastChild;
			sr(o, s);
		} else sr(t, t);
		return t;
	};
}
function cr() {
	if (T) return sr(E, null), E;
	var e = document.createDocumentFragment(), t = document.createComment(""), n = I();
	return e.append(t, n), sr(t, n), e;
}
function Z(e, t) {
	if (T) {
		var n = W;
		(!(n.f & 32768) || n.nodes.end === null) && (n.nodes.end = E), De();
		return;
	}
	e !== null && e.before(t);
}
function Q(e, t) {
	var n = t == null ? "" : typeof t == "object" ? `${t}` : t;
	n !== (e.__t ??= e.nodeValue) && (e.__t = n, e.nodeValue = `${n}`);
}
function lr(e, t) {
	return dr(e, t);
}
var ur = /* @__PURE__ */ new Map();
function dr(e, { target: t, anchor: n, props: i = {}, events: a, context: o, intro: s = !0, transformError: c }) {
	Jt();
	var l = void 0, u = pn(() => {
		var s = n ?? t.appendChild(I());
		_t(s, { pending: () => {} }, (t) => {
			Le({});
			var n = k;
			if (o && (n.c = o), a && (i.$$events = a), T && sr(t, null), l = e(t, i) || {}, T && (W.nodes.end = E, E === null || E.nodeType !== 8 || E.data !== "]")) throw Ce(), xe;
			Re();
		}, c);
		var u = /* @__PURE__ */ new Set(), d = (e) => {
			for (var n = 0; n < e.length; n++) {
				var r = e[n];
				if (!u.has(r)) {
					u.add(r);
					var i = Xn(r);
					for (let e of [t, document]) {
						var a = ur.get(e);
						a === void 0 && (a = /* @__PURE__ */ new Map(), ur.set(e, a));
						var o = a.get(r);
						o === void 0 ? (e.addEventListener(r, rr, { passive: i }), a.set(r, 1)) : a.set(r, o + 1);
					}
				}
			}
		};
		return d(r(Qn)), $n.add(d), () => {
			for (var e of u) for (let n of [t, document]) {
				var r = ur.get(n), i = r.get(e);
				--i == 0 ? (n.removeEventListener(e, rr), r.delete(e), r.size === 0 && ur.delete(n)) : r.set(e, i);
			}
			$n.delete(d), s !== n && s.parentNode?.removeChild(s);
		};
	});
	return fr.set(l, u), l;
}
var fr = /* @__PURE__ */ new WeakMap(), pr = class {
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
			if (n) wn(n), this.#r.delete(t);
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
						En(r, t), t.append(I()), this.#n.set(e, {
							effect: r,
							fragment: t
						});
					} else V(r);
					this.#r.delete(e), this.#t.delete(e);
				};
				this.#i || !n ? (this.#r.add(e), Sn(r, i, !1)) : i();
			}
		}
	};
	#o = (e) => {
		this.#e.delete(e);
		let t = Array.from(this.#e.values());
		for (let [e, n] of this.#n) t.includes(e) || (V(n.effect), this.#n.delete(e));
	};
	ensure(e, t) {
		var n = j, r = $t();
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
function $(e, t, n = !1) {
	var r;
	T && (r = E, De());
	var i = new pr(e), a = n ? x : 0;
	function o(e, t) {
		if (T) {
			var n = Ae(r);
			if (e !== parseInt(n.substring(1))) {
				var a = ke();
				D(a), i.anchor = a, Ee(!1), i.ensure(e, t), Ee(!0);
				return;
			}
		}
		i.ensure(e, t);
	}
	gn(() => {
		var e = !1;
		t((t, n = 0) => {
			e = !0, o(n, t);
		}), e || o(-1, null);
	}, a);
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/each.js
function mr(e, t) {
	return t;
}
function hr(e, t, n) {
	for (var i = [], a = t.length, o, s = t.length, c = 0; c < a; c++) {
		let n = t[c];
		Sn(n, () => {
			if (o) {
				if (o.pending.delete(n), o.done.add(n), o.pending.size === 0) {
					var t = e.outrogroups;
					gr(e, r(o.done)), t.delete(o), t.size === 0 && (e.outrogroups = null);
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
		gr(e, t, !l);
	} else o = {
		pending: new Set(t),
		done: /* @__PURE__ */ new Set()
	}, (e.outrogroups ??= /* @__PURE__ */ new Set()).add(o);
}
function gr(e, t, n = !0) {
	var r;
	if (e.pending.size > 0) {
		r = /* @__PURE__ */ new Set();
		for (let t of e.pending.values()) for (let n of t) r.add(e.items.get(n).e);
	}
	for (var i = 0; i < t.length; i++) {
		var a = t[i];
		r?.has(a) ? (a.f |= te, En(a, document.createDocumentFragment())) : V(t[i], n);
	}
}
var _r;
function vr(t, n, i, a, o, s = null) {
	var c = t, l = /* @__PURE__ */ new Map();
	if (n & 4) {
		var u = t;
		c = T ? D(/* @__PURE__ */ Yt(u)) : u.appendChild(I());
	}
	T && De();
	var d = null, f = /* @__PURE__ */ Tt(() => {
		var t = i();
		return e(t) ? t : t == null ? [] : r(t);
	}), p, m = /* @__PURE__ */ new Map(), h = !0;
	function g(e) {
		v.effect.f & 16384 || (v.pending.delete(e), v.fallback = d, br(v, p, c, n, a), d !== null && (p.length === 0 ? d.f & 33554432 ? (d.f ^= te, Sr(d, null, c)) : wn(d) : Sn(d, () => {
			d = null;
		})));
	}
	function _(e) {
		v.pending.delete(e);
	}
	var v = {
		effect: gn(() => {
			p = Y(f);
			var e = p.length;
			let t = !1;
			T && Ae(c) === "[!" != (e === 0) && (c = ke(), D(c), Ee(!1), t = !0);
			for (var r = /* @__PURE__ */ new Set(), u = j, v = $t(), y = 0; y < e; y += 1) {
				T && E.nodeType === 8 && E.data === "]" && (c = E, t = !0, Ee(!1));
				var b = p[y], x = a(b, y), S = h ? null : l.get(x);
				S ? (S.v && Lt(S.v, b), S.i && Lt(S.i, y), v && u.unskip_effect(S.e)) : (S = xr(l, h ? c : _r ??= I(), b, x, y, o, n, i), h || (S.e.f |= te), l.set(x, S)), r.add(x);
			}
			if (e === 0 && s && !d && (h ? d = B(() => s(c)) : (d = B(() => s(_r ??= I())), d.f |= te)), e > r.size && de("", "", ""), T && e > 0 && D(ke()), !h) if (m.set(u, r), v) {
				for (let [e, t] of l) r.has(e) || u.skip_effect(t.e);
				u.oncommit(g), u.ondiscard(_);
			} else g(u);
			t && Ee(!0), Y(f);
		}),
		flags: n,
		items: l,
		pending: m,
		outrogroups: null,
		fallback: d
	};
	h = !1, T && (c = E);
}
function yr(e) {
	for (; e !== null && !(e.f & 32);) e = e.next;
	return e;
}
function br(e, t, n, i, a) {
	var o = (i & 8) != 0, s = t.length, c = e.items, l = yr(e.effect.first), u, d = null, f, p = [], m = [], h, g, _, v;
	if (o) for (v = 0; v < s; v += 1) h = t[v], g = a(h, v), _ = c.get(g).e, _.f & 33554432 || (_.nodes?.a?.measure(), (f ??= /* @__PURE__ */ new Set()).add(_));
	for (v = 0; v < s; v += 1) {
		if (h = t[v], g = a(h, v), _ = c.get(g).e, e.outrogroups !== null) for (let t of e.outrogroups) t.pending.delete(_), t.done.delete(_);
		if (_.f & 33554432) if (_.f ^= te, _ === l) Sr(_, null, n);
		else {
			var y = d ? d.next : l;
			_ === e.effect.last && (e.effect.last = _.prev), _.prev && (_.prev.next = _.next), _.next && (_.next.prev = _.prev), Cr(e, d, _), Cr(e, _, y), Sr(_, y, n), d = _, p = [], m = [], l = yr(d.next);
			continue;
		}
		if (_.f & 8192 && (wn(_), o && (_.nodes?.a?.unfix(), (f ??= /* @__PURE__ */ new Set()).delete(_))), _ !== l) {
			if (u !== void 0 && u.has(_)) {
				if (p.length < m.length) {
					var b = m[0], x;
					d = b.prev;
					var S = p[0], ee = p[p.length - 1];
					for (x = 0; x < p.length; x += 1) Sr(p[x], b, n);
					for (x = 0; x < m.length; x += 1) u.delete(m[x]);
					Cr(e, S.prev, ee.next), Cr(e, d, S), Cr(e, ee, b), l = b, d = ee, --v, p = [], m = [];
				} else u.delete(_), Sr(_, l, n), Cr(e, _.prev, _.next), Cr(e, _, d === null ? e.effect.first : d.next), Cr(e, d, _), d = _;
				continue;
			}
			for (p = [], m = []; l !== null && l !== _;) (u ??= /* @__PURE__ */ new Set()).add(l), m.push(l), l = yr(l.next);
			if (l === null) continue;
		}
		_.f & 33554432 || p.push(_), d = _, l = yr(_.next);
	}
	if (e.outrogroups !== null) {
		for (let t of e.outrogroups) t.pending.size === 0 && (gr(e, r(t.done)), e.outrogroups?.delete(t));
		e.outrogroups.size === 0 && (e.outrogroups = null);
	}
	if (l !== null || u !== void 0) {
		var C = [];
		if (u !== void 0) for (_ of u) _.f & 8192 || C.push(_);
		for (; l !== null;) !(l.f & 8192) && l !== e.fallback && C.push(l), l = yr(l.next);
		var ne = C.length;
		if (ne > 0) {
			var re = i & 4 && s === 0 ? n : null;
			if (o) {
				for (v = 0; v < ne; v += 1) C[v].nodes?.a?.measure();
				for (v = 0; v < ne; v += 1) C[v].nodes?.a?.fix();
			}
			hr(e, C, re);
		}
	}
	o && He(() => {
		if (f !== void 0) for (_ of f) _.nodes?.a?.apply();
	});
}
function xr(e, t, n, r, i, a, o, s) {
	var c = o & 1 ? o & 16 ? Ft(n) : /* @__PURE__ */ It(n, !1, !1) : null, l = o & 2 ? Ft(i) : null;
	return {
		v: c,
		i: l,
		e: B(() => (a(t, c ?? n, l ?? i, s), () => {
			e.delete(r);
		}))
	};
}
function Sr(e, t, n) {
	if (e.nodes) for (var r = e.nodes.start, i = e.nodes.end, a = t && !(t.f & 33554432) ? t.nodes.start : n; r !== null;) {
		var o = /* @__PURE__ */ Xt(r);
		if (a.before(r), r === i) return;
		r = o;
	}
}
function Cr(e, t, n) {
	t === null ? e.effect.first = n : t.next = n, n === null ? e.effect.last = t : n.prev = t;
}
//#endregion
//#region node_modules/clsx/dist/clsx.mjs
function wr(e) {
	var t, n, r = "";
	if (typeof e == "string" || typeof e == "number") r += e;
	else if (typeof e == "object") if (Array.isArray(e)) {
		var i = e.length;
		for (t = 0; t < i; t++) e[t] && (n = wr(e[t])) && (r && (r += " "), r += n);
	} else for (n in e) e[n] && (r && (r += " "), r += n);
	return r;
}
function Tr() {
	for (var e, t, n = 0, r = "", i = arguments.length; n < i; n++) (e = arguments[n]) && (t = wr(e)) && (r && (r += " "), r += t);
	return r;
}
//#endregion
//#region node_modules/svelte/src/internal/shared/attributes.js
function Er(e) {
	return typeof e == "object" ? Tr(e) : e ?? "";
}
var Dr = [..." 	\n\r\f\xA0\v﻿"];
function Or(e, t, n) {
	var r = e == null ? "" : "" + e;
	if (t && (r = r ? r + " " + t : t), n) {
		for (var i of Object.keys(n)) if (n[i]) r = r ? r + " " + i : i;
		else if (r.length) for (var a = i.length, o = 0; (o = r.indexOf(i, o)) >= 0;) {
			var s = o + a;
			(o === 0 || Dr.includes(r[o - 1])) && (s === r.length || Dr.includes(r[s])) ? r = (o === 0 ? "" : r.substring(0, o)) + r.substring(s + 1) : o = s;
		}
	}
	return r === "" ? null : r;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/class.js
function kr(e, t, n, r, i, a) {
	var o = e.__className;
	if (T || o !== n || o === void 0) {
		var s = Or(n, r, a);
		(!T || s !== e.getAttribute("class")) && (s == null ? e.removeAttribute("class") : t ? e.className = s : e.setAttribute("class", s)), e.__className = n;
	} else if (a && i !== a) for (var c in a) {
		var l = !!a[c];
		(i == null || l !== !!i[c]) && e.classList.toggle(c, l);
	}
	return a;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/select.js
function Ar(t, n, r = !1) {
	if (t.multiple) {
		if (n == null) return;
		if (!e(n)) return we();
		for (var i of t.options) i.selected = n.includes(Mr(i));
		return;
	}
	for (i of t.options) if (Ut(Mr(i), n)) {
		i.selected = !0;
		return;
	}
	(!r || n !== void 0) && (t.selectedIndex = -1);
}
function jr(e) {
	var t = new MutationObserver(() => {
		Ar(e, e.__value);
	});
	t.observe(e, {
		childList: !0,
		subtree: !0,
		attributes: !0,
		attributeFilter: ["value"]
	}), un(() => {
		t.disconnect();
	});
}
function Mr(e) {
	return "__value" in e ? e.__value : e.value;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/attributes.js
var Nr = Symbol("is custom element"), Pr = Symbol("is html"), Fr = le ? "link" : "LINK", Ir = le ? "progress" : "PROGRESS";
function Lr(e) {
	if (T) {
		var t = !1, n = () => {
			if (!t) {
				if (t = !0, e.hasAttribute("value")) {
					var n = e.value;
					zr(e, "value", null), e.value = n;
				}
				if (e.hasAttribute("checked")) {
					var r = e.checked;
					zr(e, "checked", null), e.checked = r;
				}
			}
		};
		e.__on_r = n, He(n), rn();
	}
}
function Rr(e, t) {
	var n = Br(e);
	n.value === (n.value = t ?? void 0) || e.value === t && (t !== 0 || e.nodeName !== Ir) || (e.value = t ?? "");
}
function zr(e, t, n, r) {
	var i = Br(e);
	T && (i[t] = e.getAttribute(t), t === "src" || t === "srcset" || t === "href" && e.nodeName === Fr) || i[t] !== (i[t] = n) && (t === "loading" && (e[se] = n), n == null ? e.removeAttribute(t) : typeof n != "string" && Hr(e).includes(t) ? e[t] = n : e.setAttribute(t, n));
}
function Br(e) {
	return e.__attributes ??= {
		[Nr]: e.nodeName.includes("-"),
		[Pr]: e.namespaceURI === Se
	};
}
var Vr = /* @__PURE__ */ new Map();
function Hr(e) {
	var t = e.getAttribute("is") || e.nodeName, n = Vr.get(t);
	if (n) return n;
	Vr.set(t, n = []);
	for (var r, i = e, a = Element.prototype; a !== i;) {
		for (var s in r = o(i), r) r[s].set && n.push(s);
		i = l(i);
	}
	return n;
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/props.js
function Ur(e, t, n, r) {
	var i = !Fe || (n & 2) != 0, o = (n & 8) != 0, s = (n & 16) != 0, c = r, l = !0, u = () => (l && (l = !1, c = s ? Jn(r) : r), c);
	let d;
	if (o) {
		var f = ae in e || oe in e;
		d = a(e, t)?.set ?? (f && t in e ? (n) => e[t] = n : void 0);
	}
	var p, m = !1;
	o ? [p, m] = Ze(() => e[t]) : p = e[t], p === void 0 && r !== void 0 && (p = u(), d && (i && ge(t), d(p)));
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
	var _ = !1, v = (n & 1 ? Ct : Tt)(() => (_ = !1, h()));
	o && Y(v);
	var y = W;
	return (function(e, t) {
		if (arguments.length > 0) {
			let n = t ? Y(v) : i && o ? Vt(e) : e;
			return F(v, n), _ = !0, c !== void 0 && (c = n), e;
		}
		return kn && _ || y.f & 16384 ? v.v : Y(v);
	});
}
//#endregion
//#region node_modules/svelte/src/internal/disclose-version.js
typeof window < "u" && ((window.__svelte ??= {}).v ??= /* @__PURE__ */ new Set()).add("5");
//#endregion
//#region src/lib/markup.js
function Wr(e, t, n, r) {
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
function Gr(e, t) {
	return e * (1 + t / 100);
}
function Kr(e, t, n) {
	return e * Gr(t, n);
}
function qr(e, t) {
	let n = {
		materials: 0,
		labor: 0,
		equipment: 0,
		subs: 0,
		other: 0
	}, r = 0, i = 0, a = (a) => {
		let o = Wr(a.category_type, t, e.markup_overrides, e.markup_enabled), s = a.quantity * a.unit_price, c = Kr(a.quantity, a.unit_price, o);
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
function Jr(e, t) {
	let n = {
		materials: 0,
		labor: 0,
		equipment: 0,
		subs: 0,
		other: 0
	}, r = 0, i = 0;
	for (let a of e.subcategories) {
		let e = qr(a, t);
		r += e.base, i += e.withMarkup;
		for (let t of Object.keys(n)) n[t] += e.byType[t];
	}
	return {
		base: r,
		withMarkup: i,
		byType: n
	};
}
function Yr(e) {
	let t = {
		materials: 0,
		labor: 0,
		equipment: 0,
		subs: 0,
		other: 0
	}, n = 0, r = 0;
	for (let i of e.sections) {
		let a = Jr(i, e.globals);
		n += a.base, r += a.withMarkup;
		for (let e of Object.keys(t)) t[e] += a.byType[e];
	}
	return {
		base: n,
		withMarkup: r,
		byType: t
	};
}
function Xr(e) {
	return new Intl.NumberFormat("en-US", {
		style: "currency",
		currency: "USD",
		minimumFractionDigits: 2
	}).format(e);
}
function Zr(e) {
	return `${e}%`;
}
//#endregion
//#region src/lib/LineItemRow.svelte
var Qr = /* @__PURE__ */ X("<option> </option>"), $r = /* @__PURE__ */ X("<span class=\"block text-xs text-slate-400 mt-0.5 px-1\"> </span>"), ei = /* @__PURE__ */ X("<tr class=\"border-b border-slate-100 hover:bg-slate-50 text-sm group\"><td class=\"px-1 py-1 w-24\"><select></select></td><td class=\"px-1 py-1\"><input type=\"text\" class=\"w-full px-1 py-0.5 text-slate-800 bg-transparent border-0 rounded\n				focus:ring-2 focus:ring-blue-400 focus:bg-white\" placeholder=\"Item name\"/> <!></td><td class=\"px-1 py-1 w-20\"><input type=\"number\" step=\"any\" min=\"0\" class=\"w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 rounded\n				focus:ring-2 focus:ring-blue-400 focus:bg-white\"/></td><td class=\"px-1 py-1 w-16\"><input type=\"text\" class=\"w-full text-center text-slate-500 px-1 py-0.5 bg-transparent border-0 rounded\n				focus:ring-2 focus:ring-blue-400 focus:bg-white text-sm\" placeholder=\"ea\"/></td><td class=\"px-1 py-1 w-24\"><input type=\"number\" step=\"0.01\" min=\"0\" class=\"w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 rounded\n				focus:ring-2 focus:ring-blue-400 focus:bg-white\"/></td><td class=\"px-2 py-1.5 text-right font-mono text-slate-400 w-16 text-xs\"> </td><td class=\"px-2 py-1.5 text-right font-mono text-slate-500 w-24 text-xs\"> </td><td class=\"px-2 py-1.5 text-right font-mono font-medium w-28\"> </td></tr>");
function ti(e, t) {
	Le(t, !0);
	let n = Ur(t, "item", 7), r = [
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
	}, a = /* @__PURE__ */ N(() => Wr(n().category_type, t.globals, t.markupOverrides, t.markupEnabled)), o = /* @__PURE__ */ N(() => Gr(n().unit_price, Y(a))), s = /* @__PURE__ */ N(() => Kr(n().quantity, n().unit_price, Y(a))), c = /* @__PURE__ */ N(() => i[n().category_type] || i.other);
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
	var h = ei(), g = L(h), _ = L(g);
	vr(_, 21, () => r, mr, (e, t) => {
		var n = Qr(), r = L(n, !0);
		O(n);
		var i = {};
		z(() => {
			Q(r, Y(t)), i !== (i = Y(t)) && (n.value = (n.__value = Y(t)) ?? "");
		}), Z(e, n);
	}), O(_);
	var v;
	jr(_), O(g);
	var y = R(g), b = L(y);
	Lr(b);
	var x = R(b, 2), S = (e) => {
		var t = $r(), r = L(t, !0);
		O(t), z(() => Q(r, n().description)), Z(e, t);
	};
	$(x, (e) => {
		n().description && e(S);
	}), O(y);
	var ee = R(y), te = L(ee);
	Lr(te), O(ee);
	var C = R(ee), ne = L(C);
	Lr(ne), O(C);
	var re = R(C), ie = L(re);
	Lr(ie), O(re);
	var ae = R(re), oe = L(ae);
	O(ae);
	var se = R(ae), ce = L(se, !0);
	O(se);
	var le = R(se), ue = L(le, !0);
	O(le), O(h), z((e, t) => {
		kr(_, 1, `w-full text-xs px-1 py-1 rounded border-0 bg-transparent font-medium cursor-pointer
				focus:ring-2 focus:ring-blue-400 focus:bg-white ${Y(c) ?? ""}`), v !== (v = n().category_type) && (_.value = (_.__value = n().category_type) ?? "", Ar(_, n().category_type)), Rr(b, n().item_name), Rr(te, n().quantity), Rr(ne, n().unit), Rr(ie, n().unit_price), Q(oe, `${Y(a) ?? ""}%`), Q(ce, e), Q(ue, t);
	}, [() => Xr(Y(o)), () => Xr(Y(s))]), er("change", _, m), er("input", b, f), er("input", te, u), er("input", ne, p), er("input", ie, d), Z(e, h), Re();
}
tr(["change", "input"]);
//#endregion
//#region src/lib/ComponentGroupBlock.svelte
var ni = /* @__PURE__ */ X("<table class=\"w-full\"><tbody></tbody></table>"), ri = /* @__PURE__ */ X("<div class=\"ml-4 mt-2\"><div class=\"flex items-center gap-2 mb-1\"><span class=\"text-xs font-semibold text-slate-500 uppercase tracking-wide\"> </span> <span class=\"text-xs text-slate-400\"> </span></div> <!></div>");
function ii(e, t) {
	Le(t, !0);
	var n = ri(), r = L(n), i = L(r), a = L(i, !0);
	O(i);
	var o = R(i, 2), s = L(o);
	O(o), O(r);
	var c = R(r, 2), l = (e) => {
		var n = ni(), r = L(n);
		vr(r, 21, () => t.group.line_items, (e) => e.id, (e, n) => {
			ti(e, {
				get item() {
					return Y(n);
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
				}
			});
		}), O(r), O(n), Z(e, n);
	};
	$(c, (e) => {
		t.group.line_items.length > 0 && e(l);
	}), O(n), z(() => {
		Q(a, t.group.name), Q(s, `(${t.group.line_items.length ?? ""})`);
	}), Z(e, n), Re();
}
//#endregion
//#region src/lib/SubcategoryBlock.svelte
var ai = /* @__PURE__ */ X("<span class=\"text-xs px-1.5 py-0.5 rounded bg-amber-50 text-amber-600 font-medium\">overrides</span>"), oi = /* @__PURE__ */ X("<span class=\"text-xs px-1.5 py-0.5 rounded bg-green-50 text-green-600 font-medium\"> </span>"), si = /* @__PURE__ */ X("<span> </span>"), ci = /* @__PURE__ */ X("<div class=\"flex items-center gap-3 py-1.5 px-2 bg-amber-50 rounded text-xs mb-2\"><span class=\"font-medium text-amber-700\">Markup:</span> <!></div>"), li = /* @__PURE__ */ X("<table class=\"w-full\"><thead><tr class=\"text-xs text-slate-400 uppercase tracking-wide border-b border-slate-100\"><th class=\"px-1 py-1 text-left w-24\">Type</th><th class=\"px-1 py-1 text-left\">Name</th><th class=\"px-1 py-1 text-right w-20\">Qty</th><th class=\"px-1 py-1 text-center w-16\">Unit</th><th class=\"px-1 py-1 text-right w-24\">Price</th><th class=\"px-2 py-1 text-right w-16\">Markup</th><th class=\"px-2 py-1 text-right w-24\">w/ Markup</th><th class=\"px-2 py-1 text-right w-28\">Total</th></tr></thead><tbody></tbody></table>"), ui = /* @__PURE__ */ X("<span class=\"text-green-600\"> </span>"), di = /* @__PURE__ */ X("<div class=\"px-4 pb-3\"><!> <!> <!> <div class=\"flex justify-between items-center mt-2 pt-2 border-t border-slate-100 text-sm\"><span class=\"text-slate-500\"> <!></span> <span class=\"font-mono font-semibold text-slate-700\"> </span></div></div>"), fi = /* @__PURE__ */ X("<div class=\"border-t border-slate-200\"><button class=\"w-full flex items-center justify-between px-4 py-2 hover:bg-slate-50 transition-colors text-left\"><div class=\"flex items-center gap-2\"><svg fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> <span class=\"font-medium text-slate-600 text-sm\"> </span> <span class=\"text-xs text-slate-400\"> </span> <!> <!></div> <span class=\"font-mono text-sm font-semibold text-slate-700\"> </span></button> <!></div>");
function pi(e, t) {
	Le(t, !0);
	let n = /* @__PURE__ */ P(Vt(Ur(t, "collapsed", 3, !1)())), r = /* @__PURE__ */ N(() => qr(t.subcat, t.globals)), i = /* @__PURE__ */ N(() => [
		"materials",
		"labor",
		"equipment",
		"subs",
		"other"
	].map((e) => ({
		type: e,
		value: Wr(e, t.globals, t.subcat.markup_overrides, t.subcat.markup_enabled),
		isOverride: t.subcat.markup_overrides[e] != null,
		isDisabled: !t.subcat.markup_enabled[e]
	}))), a = /* @__PURE__ */ N(() => Y(i).some((e) => e.isOverride || e.isDisabled)), o = /* @__PURE__ */ N(() => t.subcat.line_items.length + t.subcat.component_groups.reduce((e, t) => e + t.line_items.length, 0));
	var s = fi(), c = L(s), l = L(c), u = L(l), d = R(u, 2), f = L(d, !0);
	O(d);
	var p = R(d, 2), m = L(p);
	O(p);
	var h = R(p, 2), g = (e) => {
		Z(e, ai());
	};
	$(h, (e) => {
		Y(a) && e(g);
	});
	var _ = R(h, 2), v = (e) => {
		var n = oi(), r = L(n);
		O(n), z((e) => Q(r, `+${e ?? ""} lump sum`), [() => Xr(t.subcat.lump_sum)]), Z(e, n);
	};
	$(_, (e) => {
		t.subcat.lump_sum > 0 && e(v);
	}), O(l);
	var y = R(l, 2), b = L(y, !0);
	O(y), O(c);
	var x = R(c, 2), S = (e) => {
		var n = di(), s = L(n), c = (e) => {
			var t = ci();
			vr(R(L(t), 2), 17, () => Y(i), mr, (e, t) => {
				var n = si(), r = L(n);
				O(n), z((e) => {
					kr(n, 1, Er(Y(t).isDisabled ? "text-slate-400 line-through" : Y(t).isOverride ? "text-amber-700 font-medium" : "text-slate-500")), Q(r, `${Y(t).type ?? ""} ${e ?? ""}`);
				}, [() => Zr(Y(t).value)]), Z(e, n);
			}), O(t), Z(e, t);
		};
		$(s, (e) => {
			Y(a) && e(c);
		});
		var l = R(s, 2), u = (e) => {
			var n = li(), r = R(L(n));
			vr(r, 21, () => t.subcat.line_items, (e) => e.id, (e, n) => {
				ti(e, {
					get item() {
						return Y(n);
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
					}
				});
			}), O(r), O(n), Z(e, n);
		};
		$(l, (e) => {
			Y(o) > 0 && e(u);
		});
		var d = R(l, 2);
		vr(d, 17, () => t.subcat.component_groups, (e) => e.id, (e, n) => {
			ii(e, {
				get group() {
					return Y(n);
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
				}
			});
		});
		var f = R(d, 2), p = L(f), m = L(p), h = R(m), g = (e) => {
			var n = ui(), r = L(n);
			O(n), z((e) => Q(r, `+ ${e ?? ""} lump sum`), [() => Xr(t.subcat.lump_sum)]), Z(e, n);
		};
		$(h, (e) => {
			t.subcat.lump_sum > 0 && e(g);
		}), O(p);
		var _ = R(p, 2), v = L(_, !0);
		O(_), O(f), O(n), z((e, t) => {
			Q(m, `Subtotal: ${e ?? ""} `), Q(v, t);
		}, [() => Xr(Y(r).base), () => Xr(Y(r).withMarkup)]), Z(e, n);
	};
	$(x, (e) => {
		Y(n) || e(S);
	}), O(s), z((e) => {
		kr(u, 0, `w-4 h-4 text-slate-400 transition-transform ${Y(n) ? "" : "rotate-90"}`), Q(f, t.subcat.name), Q(m, `${Y(o) ?? ""} item${Y(o) === 1 ? "" : "s"}`), Q(b, e);
	}, [() => Xr(Y(r).withMarkup)]), er("click", c, () => F(n, !Y(n))), Z(e, s), Re();
}
tr(["click"]);
//#endregion
//#region src/lib/SectionBlock.svelte
var mi = /* @__PURE__ */ X("<div class=\"px-4 py-8 text-center text-slate-400 text-sm\">No subcategories yet</div>"), hi = /* @__PURE__ */ X("<div class=\"mb-4 border border-slate-200 rounded-lg overflow-hidden bg-white\"><button class=\"w-full flex items-center justify-between px-4 py-3 bg-slate-800 text-white hover:bg-slate-700 transition-colors text-left\"><div class=\"flex items-center gap-3\"><svg fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> <span class=\"font-semibold\"> </span> <span class=\"text-xs text-slate-400\"> </span></div> <span class=\"font-mono font-semibold\"> </span></button> <!></div>");
function gi(e, t) {
	Le(t, !0);
	let n = /* @__PURE__ */ P(Vt(Ur(t, "collapsed", 3, !1)())), r = /* @__PURE__ */ N(() => Jr(t.section, t.globals)), i = /* @__PURE__ */ N(() => t.section.subcategories.reduce((e, t) => e + t.line_items.length + t.component_groups.reduce((e, t) => e + t.line_items.length, 0), 0));
	var a = hi(), o = L(a), s = L(o), c = L(s), l = R(c, 2), u = L(l, !0);
	O(l);
	var d = R(l, 2), f = L(d);
	O(d), O(s);
	var p = R(s, 2), m = L(p, !0);
	O(p), O(o);
	var h = R(o, 2), g = (e) => {
		var n = cr(), r = Zt(n), i = (e) => {
			Z(e, mi());
		}, a = (e) => {
			var n = cr();
			vr(Zt(n), 17, () => t.section.subcategories, (e) => e.id, (e, n) => {
				pi(e, {
					get subcat() {
						return Y(n);
					},
					get globals() {
						return t.globals;
					},
					get onchange() {
						return t.onchange;
					}
				});
			}), Z(e, n);
		};
		$(r, (e) => {
			t.section.subcategories.length === 0 ? e(i) : e(a, -1);
		}), Z(e, n);
	};
	$(h, (e) => {
		Y(n) || e(g);
	}), O(a), z((e) => {
		kr(c, 0, `w-4 h-4 text-slate-400 transition-transform ${Y(n) ? "" : "rotate-90"}`), Q(u, t.section.name), Q(f, `${t.section.subcategories.length ?? ""} subcategor${t.section.subcategories.length === 1 ? "y" : "ies"}
				· ${Y(i) ?? ""} item${Y(i) === 1 ? "" : "s"}`), Q(m, e);
	}, [() => Xr(Y(r).withMarkup)]), er("click", o, () => F(n, !Y(n))), Z(e, a), Re();
}
tr(["click"]);
//#endregion
//#region src/lib/FooterSummary.svelte
var _i = /* @__PURE__ */ X("<div class=\"flex flex-col\"><span class=\"text-xs text-slate-400 uppercase tracking-wide\"> </span> <span class=\"font-mono\"> </span></div>"), vi = /* @__PURE__ */ X("<span class=\"text-slate-500 text-xs\">No items yet</span>"), yi = /* @__PURE__ */ X("<div class=\"fixed bottom-0 left-0 right-0 bg-slate-900 text-white border-t border-slate-700 z-20\"><div class=\"flex items-center justify-between px-4 py-3\"><div class=\"flex items-center gap-6 text-sm\"><div class=\"flex flex-col\"><span class=\"text-xs text-slate-400 uppercase tracking-wide\">Base Cost</span> <span class=\"font-mono\"> </span></div> <div class=\"w-px h-8 bg-slate-700\"></div> <!> <!></div> <div class=\"flex flex-col items-end\"><span class=\"text-xs text-slate-400 uppercase tracking-wide\">Total</span> <span class=\"font-mono text-lg font-bold\"> </span></div></div></div>");
function bi(e, t) {
	Le(t, !0);
	let n = /* @__PURE__ */ N(() => Yr(t.estimate)), r = /* @__PURE__ */ N(() => Object.entries(Y(n).byType).filter(([, e]) => e > 0).map(([e, t]) => ({
		type: e,
		value: t
	}))), i = {
		materials: "Materials",
		labor: "Labor",
		equipment: "Equipment",
		subs: "Subs",
		other: "Other"
	};
	var a = yi(), o = L(a), s = L(o), c = L(s), l = R(L(c), 2), u = L(l, !0);
	O(l), O(c);
	var d = R(c, 4);
	vr(d, 17, () => Y(r), mr, (e, t) => {
		let n = () => Y(t).type, r = () => Y(t).value;
		var a = _i(), o = L(a), s = L(o, !0);
		O(o);
		var c = R(o, 2), l = L(c, !0);
		O(c), O(a), z((e) => {
			Q(s, i[n()]), Q(l, e);
		}, [() => Xr(r())]), Z(e, a);
	});
	var f = R(d, 2), p = (e) => {
		Z(e, vi());
	};
	$(f, (e) => {
		Y(r).length === 0 && e(p);
	}), O(s);
	var m = R(s, 2), h = R(L(m), 2), g = L(h, !0);
	O(h), O(m), O(o), O(a), z((e, t) => {
		Q(u, e), Q(g, t);
	}, [() => Xr(Y(n).base), () => Xr(Y(n).withMarkup)]), Z(e, a), Re();
}
//#endregion
//#region src/lib/SaveStatus.svelte
var xi = /* @__PURE__ */ X("<div><span></span> </div>");
function Si(e, t) {
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
	var a = xi(), o = L(a), s = R(o);
	O(a), z(() => {
		kr(a, 1, `flex items-center gap-1.5 text-xs ${Y(r) ?? ""}`), kr(o, 1, `w-2 h-2 rounded-full ${Y(i) ?? ""}`), Q(s, ` ${Y(n) ?? ""}`);
	}), Z(e, a), Re();
}
//#endregion
//#region src/lib/autosave.svelte.js
function Ci(e, t = 2e3) {
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
			return Y(n);
		},
		get savedAt() {
			return Y(r);
		}
	};
}
//#endregion
//#region src/EstimateBuilder.svelte
var wi = /* @__PURE__ */ X("<div class=\"flex items-center justify-center h-64\"><div class=\"text-slate-500\">Loading estimate...</div></div>"), Ti = /* @__PURE__ */ X("<div class=\"flex items-center justify-center h-64\"><div class=\"text-red-500\"> </div></div>"), Ei = /* @__PURE__ */ X("<div class=\"text-center py-16 text-slate-400\"><svg class=\"w-12 h-12 mx-auto mb-3 text-slate-300\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"1.5\" d=\"M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z\"></path></svg> <p class=\"text-lg font-medium\">No sections yet</p> <p class=\"text-sm mt-1\">Add a section to start building your estimate.</p></div>"), Di = /* @__PURE__ */ X("<div class=\"estimate-builder pb-16\"><div class=\"sticky top-0 z-10 bg-slate-900 text-white px-4 py-3 flex items-center justify-between border-b border-slate-700\"><div><h1 class=\"text-lg font-semibold\"> </h1> <span class=\"text-xs text-slate-400\">Estimate Builder</span></div> <div class=\"flex items-center gap-3\"><!> <span class=\"text-xs px-2 py-1 rounded bg-slate-700 text-slate-300 uppercase tracking-wide\"> </span></div></div> <div class=\"bg-slate-50 border-b border-slate-200 px-4 py-2\"><div class=\"flex items-center gap-4 text-sm\"><span class=\"font-medium text-slate-600\">Global Markup:</span> <span class=\"px-2 py-0.5 rounded bg-blue-50 text-blue-700 font-mono text-xs\"> </span> <span class=\"px-2 py-0.5 rounded bg-amber-50 text-amber-700 font-mono text-xs\"> </span> <span class=\"px-2 py-0.5 rounded bg-purple-50 text-purple-700 font-mono text-xs\"> </span> <span class=\"px-2 py-0.5 rounded bg-green-50 text-green-700 font-mono text-xs\"> </span> <span class=\"px-2 py-0.5 rounded bg-slate-100 text-slate-600 font-mono text-xs\"> </span></div></div> <div class=\"p-4\"><!></div> <!></div>");
function Oi(e, t) {
	Le(t, !0);
	let n = /* @__PURE__ */ P(null), r = /* @__PURE__ */ P(null), i = /* @__PURE__ */ P(!0), a = Ci(t.projectId);
	async function o() {
		try {
			let e = await fetch(`/api/estimate/${t.projectId}`);
			if (!e.ok) throw Error(`Failed to load estimate: ${e.status}`);
			F(n, await e.json(), !0), a.register(() => Y(n));
		} catch (e) {
			F(r, e.message, !0);
		} finally {
			F(i, !1);
		}
	}
	dn(() => {
		t.projectId && o();
	}), dn(() => {
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
	var c = cr(), l = Zt(c), u = (e) => {
		Z(e, wi());
	}, d = (e) => {
		var t = Ti(), n = L(t), i = L(n, !0);
		O(n), O(t), z(() => Q(i, Y(r))), Z(e, t);
	}, f = (e) => {
		var t = Di(), r = L(t), i = L(r), o = L(i), c = L(o, !0);
		O(o), Oe(2), O(i);
		var l = R(i, 2), u = L(l);
		Si(u, {
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
			Z(e, Ei());
		}, ie = (e) => {
			var t = cr();
			vr(Zt(t), 17, () => Y(n).sections, (e) => e.id, (e, t) => {
				gi(e, {
					get section() {
						return Y(t);
					},
					get globals() {
						return Y(n).globals;
					},
					onchange: s
				});
			}), Z(e, t);
		};
		$(ne, (e) => {
			Y(n).sections.length === 0 ? e(re) : e(ie, -1);
		}), O(C), bi(R(C, 2), { get estimate() {
			return Y(n);
		} }), O(t), z((e, t, r, i, a) => {
			Q(c, Y(n).project.name), Q(f, Y(n).project.status), Q(g, `Materials ${e ?? ""}`), Q(v, `Labor ${t ?? ""}`), Q(b, `Equipment ${r ?? ""}`), Q(S, `Subs ${i ?? ""}`), Q(te, `Other ${a ?? ""}`);
		}, [
			() => Zr(Y(n).globals.materials_markup),
			() => Zr(Y(n).globals.labor_markup),
			() => Zr(Y(n).globals.equipment_markup),
			() => Zr(Y(n).globals.subs_markup),
			() => Zr(Y(n).globals.other_markup)
		]), Z(e, t);
	};
	$(l, (e) => {
		Y(i) ? e(u) : Y(r) ? e(d, 1) : Y(n) && e(f, 2);
	}), Z(e, c), Re();
}
//#endregion
//#region src/main.js
var ki = document.getElementById("estimate-root");
if (ki) {
	let e = ki.dataset.projectId;
	lr(Oi, {
		target: ki,
		props: { projectId: e }
	});
}
//#endregion
