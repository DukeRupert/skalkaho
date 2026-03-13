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
var p = 1024, m = 2048, h = 4096, g = 8192, _ = 16384, v = 32768, y = 1 << 25, b = 65536, x = 1 << 19, S = 1 << 20, ee = 1 << 25, C = 65536, te = 1 << 21, ne = 1 << 22, re = 1 << 23, ie = Symbol("$state"), ae = Symbol("legacy props"), oe = new class extends Error {
	name = "StaleReactionError";
	message = "The reaction that called `getAbortSignal()` was re-run or destroyed";
}();
globalThis.document?.contentType;
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
var ve = {}, w = Symbol();
function ye(e) {
	console.warn("https://svelte.dev/e/hydration_mismatch");
}
function be() {
	console.warn("https://svelte.dev/e/svelte_boundary_reset_noop");
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/hydration.js
var T = !1;
function xe(e) {
	T = e;
}
var E;
function D(e) {
	if (e === null) throw ye(), ve;
	return E = e;
}
function Se() {
	return D(/* @__PURE__ */ Ut(E));
}
function O(e) {
	if (T) {
		if (/* @__PURE__ */ Ut(E) !== null) throw ye(), ve;
		E = e;
	}
}
function Ce(e = 1) {
	if (T) {
		for (var t = e, n = E; t--;) n = /* @__PURE__ */ Ut(n);
		E = n;
	}
}
function we(e = !0) {
	for (var t = 0, n = E;;) {
		if (n.nodeType === 8) {
			var r = n.data;
			if (r === "]") {
				if (t === 0) return n;
				--t;
			} else (r === "[" || r === "[!" || r[0] === "[" && !isNaN(Number(r.slice(1)))) && (t += 1);
		}
		var i = /* @__PURE__ */ Ut(n);
		e && n.remove(), n = i;
	}
}
function Te(e) {
	if (!e || e.nodeType !== 8) throw ye(), ve;
	return e.data;
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/equality.js
function Ee(e) {
	return e === this.v;
}
function De(e, t) {
	return e == e ? e !== t || typeof e == "object" && !!e || typeof e == "function" : t == t;
}
function Oe(e) {
	return !De(e, this.v);
}
//#endregion
//#region node_modules/svelte/src/internal/flags/index.js
var ke = !1, Ae = !1, k = null;
function je(e) {
	k = e;
}
function Me(e, t = !1, n) {
	k = {
		p: k,
		i: !1,
		c: null,
		e: null,
		s: e,
		x: null,
		r: W,
		l: Ae && !t ? {
			s: null,
			u: null,
			$: []
		} : null
	};
}
function Ne(e) {
	var t = k, n = t.e;
	if (n !== null) {
		t.e = null;
		for (var r of n) nn(r);
	}
	return e !== void 0 && (t.x = e), t.i = !0, k = t.p, e ?? {};
}
function Pe() {
	return !Ae || k !== null && k.l === null;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/task.js
var Fe = [];
function Ie() {
	var e = Fe;
	Fe = [], d(e);
}
function Le(e) {
	if (Fe.length === 0 && !Ye) {
		var t = Fe;
		queueMicrotask(() => {
			t === Fe && Ie();
		});
	}
	Fe.push(e);
}
function Re(e) {
	var t = W;
	if (t === null) return H.f |= re, e;
	if (!(t.f & 32768) && !(t.f & 4)) throw e;
	ze(e, t);
}
function ze(e, t) {
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
var Be = ~(m | h | p);
function A(e, t) {
	e.f = e.f & Be | t;
}
function Ve(e) {
	e.f & 512 || e.deps === null ? A(e, p) : A(e, h);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/utils.js
function He(e) {
	if (e !== null) for (let t of e) !(t.f & 2) || !(t.f & 65536) || (t.f ^= C, He(t.deps));
}
function Ue(e, t, n) {
	e.f & 2048 ? t.add(e) : e.f & 4096 && n.add(e), He(e.deps), A(e, p);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/store.js
var We = !1, Ge = !1;
function Ke(e) {
	var t = Ge;
	try {
		return Ge = !1, [e(), Ge];
	} finally {
		Ge = t;
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/batch.js
var qe = /* @__PURE__ */ new Set(), j = null, M = null, Je = null, Ye = !1, Xe = !1, Ze = null, Qe = null, $e = 0, et = 1, tt = class e {
	id = et++;
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
		$e++ > 1e3 && nt();
		let t = this.#a;
		this.#a = [], this.apply();
		var n = Ze = [], r = [], i = Qe = [];
		for (let e of t) try {
			this.#f(e, n, r);
		} catch (t) {
			throw lt(e), t;
		}
		if (j = null, i.length > 0) {
			var a = e.ensure();
			for (let e of i) a.schedule(e);
		}
		if (Ze = null, Qe = null, this.#u()) {
			this.#p(r), this.#p(n);
			for (let [e, t] of this.#c) ct(e, t);
		} else {
			this.#n === 0 && qe.delete(this), this.#o.clear(), this.#s.clear();
			for (let e of this.#e) e(this);
			this.#e.clear(), it(r), it(n), this.#i?.resolve();
		}
		var o = j;
		if (this.#a.length > 0) {
			let e = o ??= this;
			e.#a.push(...this.#a.filter((t) => !e.#a.includes(t)));
		}
		o !== null && (qe.add(o), o.#d()), qe.has(this) || this.#m();
	}
	#f(e, t, n) {
		e.f ^= p;
		for (var r = e.first; r !== null;) {
			var i = r.f, a = (i & 96) != 0;
			if (!(a && i & 1024 || i & 8192 || this.#c.has(r)) && r.fn !== null) {
				a ? r.f ^= p : i & 4 ? t.push(r) : ke && i & 16777224 ? n.push(r) : jn(r) && (i & 16 && this.#s.add(r), In(r));
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
		for (var t = 0; t < e.length; t += 1) Ue(e[t], this.#o, this.#s);
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
			if (Xe = !0, j = this, !this.#u()) {
				for (let e of this.#o) this.#s.delete(e), A(e, m), this.schedule(e);
				for (let e of this.#s) A(e, h), this.schedule(e);
			}
			this.#d();
		} finally {
			$e = 0, Je = null, Ze = null, Qe = null, Xe = !1, j = null, M = null, Ot.clear();
		}
	}
	discard() {
		for (let e of this.#t) e(this);
		this.#t.clear();
	}
	#m() {
		for (let s of qe) {
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
					for (var a of t) at(a, n, r, i);
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
		--this.#n, e && --this.#r, !(this.#l || t) && (this.#l = !0, Le(() => {
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
			Xe || (qe.add(j), Ye || Le(() => {
				j === t && t.flush();
			}));
		}
		return j;
	}
	apply() {
		if (!ke || !this.is_fork && qe.size === 1) {
			M = null;
			return;
		}
		M = new Map(this.current);
		for (let e of qe) if (e !== this) for (let [t, n] of e.previous) M.has(t) || M.set(t, n);
	}
	schedule(e) {
		if (Je = e, e.b?.is_pending && e.f & 16777228 && !(e.f & 32768)) {
			e.b.defer_effect(e);
			return;
		}
		for (var t = e; t.parent !== null;) {
			t = t.parent;
			var n = t.f;
			if (Ze !== null && t === W && (ke || (H === null || !(H.f & 2)) && !We)) return;
			if (n & 96) {
				if (!(n & 1024)) return;
				t.f ^= p;
			}
		}
		this.#a.push(t);
	}
};
function nt() {
	try {
		fe();
	} catch (e) {
		ze(e, Je);
	}
}
var rt = null;
function it(e) {
	var t = e.length;
	if (t !== 0) {
		for (var n = 0; n < t;) {
			var r = e[n++];
			if (!(r.f & 24576) && jn(r) && (rt = /* @__PURE__ */ new Set(), In(r), r.deps === null && r.first === null && r.nodes === null && r.teardown === null && r.ac === null && fn(r), rt?.size > 0)) {
				Ot.clear();
				for (let e of rt) {
					if (e.f & 24576) continue;
					let t = [e], n = e.parent;
					for (; n !== null;) rt.has(n) && (rt.delete(n), t.push(n)), n = n.parent;
					for (let e = t.length - 1; e >= 0; e--) {
						let n = t[e];
						n.f & 24576 || In(n);
					}
				}
				rt.clear();
			}
		}
		rt = null;
	}
}
function at(e, t, n, r) {
	if (!n.has(e) && (n.add(e), e.reactions !== null)) for (let i of e.reactions) {
		let e = i.f;
		e & 2 ? at(i, t, n, r) : e & 4194320 && !(e & 2048) && ot(i, t, r) && (A(i, m), st(i));
	}
}
function ot(e, t, r) {
	let i = r.get(e);
	if (i !== void 0) return i;
	if (e.deps !== null) for (let i of e.deps) {
		if (n.call(t, i)) return !0;
		if (i.f & 2 && ot(i, t, r)) return r.set(i, !0), !0;
	}
	return r.set(e, !1), !1;
}
function st(e) {
	j.schedule(e);
}
function ct(e, t) {
	if (!(e.f & 32 && e.f & 1024)) {
		e.f & 2048 ? t.d.push(e) : e.f & 4096 && t.m.push(e), A(e, p);
		for (var n = e.first; n !== null;) ct(n, t), n = n.next;
	}
}
function lt(e) {
	A(e, p);
	for (var t = e.first; t !== null;) lt(t), t = t.next;
}
//#endregion
//#region node_modules/svelte/src/reactivity/create-subscriber.js
function ut(e) {
	let t = 0, n = At(0), r;
	return () => {
		$t() && (Y(n), on(() => (t === 0 && (r = zn(() => e(() => Pt(n)))), t += 1, () => {
			Le(() => {
				--t, t === 0 && (r?.(), r = void 0, Pt(n));
			});
		})));
	};
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/boundary.js
var dt = b | x;
function ft(e, t, n, r) {
	new pt(e, t, n, r);
}
var pt = class {
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
	#h = ut(() => (this.#m = At(this.#l), () => {
		this.#m = null;
	}));
	constructor(e, t, n, r) {
		this.#e = e, this.#n = t, this.#r = (e) => {
			var t = W;
			t.b = this, t.f |= 128, n(e);
		}, this.parent = W.b, this.transform_error = r ?? this.parent?.transform_error ?? ((e) => e), this.#i = sn(() => {
			if (T) {
				let e = this.#t;
				Se();
				let t = e.data === "[!";
				if (e.data.startsWith("[?")) {
					let t = JSON.parse(e.data.slice(2));
					this.#_(t);
				} else t ? this.#v() : this.#g();
			} else this.#y();
		}, dt), T && (this.#e = E);
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
		e && (this.is_pending = !0, this.#o = B(() => e(this.#e)), Le(() => {
			var e = this.#c = document.createDocumentFragment(), t = I();
			e.append(t), this.#a = this.#x(() => B(() => this.#r(t))), this.#u === 0 && (this.#e.before(e), this.#c = null, pn(this.#o, () => {
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
				_n(this.#a, e);
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
		Ue(e, this.#f, this.#p);
	}
	is_rendered() {
		return !this.is_pending && (!this.parent || this.parent.is_rendered());
	}
	has_pending_snippet() {
		return !!this.#n.pending;
	}
	#x(e) {
		var t = W, n = H, r = k;
		Cn(this.#i), U(this.#i), je(this.#i.ctx);
		try {
			return tt.ensure(), e();
		} catch (e) {
			return Re(e), null;
		} finally {
			Cn(t), U(n), je(r);
		}
	}
	#S(e, t) {
		if (!this.has_pending_snippet()) {
			this.parent && this.parent.#S(e, t);
			return;
		}
		this.#u += e, this.#u === 0 && (this.#b(t), this.#o && pn(this.#o, () => {
			this.#o = null;
		}), this.#c &&= (this.#e.before(this.#c), null));
	}
	update_pending_count(e, t) {
		this.#S(e, t), this.#l += e, !(!this.#m || this.#d) && (this.#d = !0, Le(() => {
			this.#d = !1, this.#m && Mt(this.#m, this.#l);
		}));
	}
	get_effect_pending() {
		return this.#h(), Y(this.#m);
	}
	error(e) {
		var t = this.#n.onerror;
		let n = this.#n.failed;
		if (!t && !n) throw e;
		this.#a &&= (V(this.#a), null), this.#o &&= (V(this.#o), null), this.#s &&= (V(this.#s), null), T && (D(this.#t), Ce(), D(we()));
		var r = !1, i = !1;
		let a = () => {
			if (r) {
				be();
				return;
			}
			r = !0, i && _e(), this.#s !== null && pn(this.#s, () => {
				this.#s = null;
			}), this.#x(() => {
				this.#y();
			});
		}, o = (e) => {
			try {
				i = !0, t?.(e, a), i = !1;
			} catch (e) {
				ze(e, this.#i && this.#i.parent);
			}
			n && (this.#s = this.#x(() => {
				try {
					return B(() => {
						var t = W;
						t.b = this, t.f |= 128, n(this.#e, () => e, () => a);
					});
				} catch (e) {
					return ze(e, this.#i.parent), null;
				}
			}));
		};
		Le(() => {
			var t;
			try {
				t = this.transform_error(e);
			} catch (e) {
				ze(e, this.#i && this.#i.parent);
				return;
			}
			typeof t == "object" && t && typeof t.then == "function" ? t.then(o, (e) => ze(e, this.#i && this.#i.parent)) : o(t);
		});
	}
};
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/async.js
function mt(e, t, n, r) {
	let i = Pe() ? vt : bt;
	var a = e.filter((e) => !e.settled);
	if (n.length === 0 && a.length === 0) {
		r(t.map(i));
		return;
	}
	var o = W, s = ht(), c = a.length === 1 ? a[0].promise : a.length > 1 ? Promise.all(a.map((e) => e.promise)) : null;
	function l(e) {
		s();
		try {
			r(e);
		} catch (e) {
			o.f & 16384 || ze(e, o);
		}
		gt();
	}
	if (n.length === 0) {
		c.then(() => l(t.map(i)));
		return;
	}
	var u = _t();
	function d() {
		Promise.all(n.map((e) => /* @__PURE__ */ yt(e))).then((e) => l([...t.map(i), ...e])).catch((e) => ze(e, o)).finally(() => u());
	}
	c ? c.then(() => {
		s(), d(), gt();
	}) : d();
}
function ht() {
	var e = W, t = H, n = k, r = j;
	return function(i = !0) {
		Cn(e), U(t), je(n), i && !(e.f & 16384) && (r?.activate(), r?.apply());
	};
}
function gt(e = !0) {
	Cn(null), U(null), je(null), e && j?.deactivate();
}
function _t() {
	var e = W.b, t = j, n = e.is_rendered();
	return e.update_pending_count(1, t), t.increment(n), (r = !1) => {
		e.update_pending_count(-1, t), t.decrement(n, r);
	};
}
/* @__NO_SIDE_EFFECTS__ */
function vt(e) {
	var t = 2 | m, n = H !== null && H.f & 2 ? H : null;
	return W !== null && (W.f |= x), {
		ctx: k,
		deps: null,
		effects: null,
		equals: Ee,
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
function yt(e, t, n) {
	let r = W;
	r === null && se();
	var i = void 0, a = At(w), o = !H, s = /* @__PURE__ */ new Map();
	return an(() => {
		var t = W, n = f();
		i = n.promise;
		try {
			Promise.resolve(e()).then(n.resolve, n.reject).finally(gt);
		} catch (e) {
			n.reject(e), gt();
		}
		var c = j;
		if (o) {
			if (t.f & 32768) var l = _t();
			if (r.b.is_rendered()) s.get(c)?.reject(oe), s.delete(c);
			else {
				for (let e of s.values()) e.reject(oe);
				s.clear();
			}
			s.set(c, n);
		}
		let u = (e, n = void 0) => {
			if (l && l(n === oe), !(n === oe || t.f & 16384)) {
				if (c.activate(), n) a.f |= re, Mt(a, n);
				else {
					a.f & 8388608 && (a.f ^= re), Mt(a, e);
					for (let [e, t] of s) {
						if (s.delete(e), e === c) break;
						t.reject(oe);
					}
				}
				c.deactivate();
			}
		};
		n.promise.then(u, (e) => u(null, e || "unknown"));
	}), en(() => {
		for (let e of s.values()) e.reject(oe);
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
	let t = /* @__PURE__ */ vt(e);
	return ke || wn(t), t;
}
/* @__NO_SIDE_EFFECTS__ */
function bt(e) {
	let t = /* @__PURE__ */ vt(e);
	return t.equals = Oe, t;
}
function xt(e) {
	var t = e.effects;
	if (t !== null) {
		e.effects = null;
		for (var n = 0; n < t.length; n += 1) V(t[n]);
	}
}
function St(e) {
	for (var t = e.parent; t !== null;) {
		if (!(t.f & 2)) return t.f & 16384 ? null : t;
		t = t.parent;
	}
	return null;
}
function Ct(e) {
	var t, n = W;
	Cn(St(e));
	try {
		e.f &= ~C, xt(e), t = Nn(e);
	} finally {
		Cn(n);
	}
	return t;
}
function wt(e) {
	var t = Ct(e);
	if (!e.equals(t) && (e.wv = An(), (!j?.is_fork || e.deps === null) && (e.v = t, e.deps === null))) {
		A(e, p);
		return;
	}
	bn || (M === null ? Ve(e) : ($t() || j?.is_fork) && M.set(e, t));
}
function Tt(e) {
	if (e.effects !== null) for (let t of e.effects) (t.teardown || t.ac) && (t.teardown?.(), t.ac?.abort(oe), t.teardown = u, t.ac = null, Fn(t, 0), ln(t));
}
function Et(e) {
	if (e.effects !== null) for (let t of e.effects) t.teardown && In(t);
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/sources.js
var Dt = /* @__PURE__ */ new Set(), Ot = /* @__PURE__ */ new Map(), kt = !1;
function At(e, t) {
	return {
		f: 0,
		v: e,
		reactions: null,
		equals: Ee,
		rv: 0,
		wv: 0
	};
}
/* @__NO_SIDE_EFFECTS__ */
function P(e, t) {
	let n = At(e, t);
	return wn(n), n;
}
/* @__NO_SIDE_EFFECTS__ */
function jt(e, t = !1, n = !0) {
	let r = At(e);
	return t || (r.equals = Oe), Ae && n && k !== null && k.l !== null && (k.l.s ??= []).push(r), r;
}
function F(e, t, r = !1) {
	return H !== null && (!Sn || H.f & 131072) && Pe() && H.f & 4325394 && (G === null || !n.call(G, e)) && ge(), Mt(e, r ? It(t) : t, Qe);
}
function Mt(e, t, n = null) {
	if (!e.equals(t)) {
		var r = e.v;
		bn ? Ot.set(e, t) : Ot.set(e, r), e.v = t;
		var i = tt.ensure();
		if (i.capture(e, r), e.f & 2) {
			let t = e;
			e.f & 2048 && Ct(t), Ve(t);
		}
		e.wv = An(), Ft(e, m, n), Pe() && W !== null && W.f & 1024 && !(W.f & 96) && (J === null ? Tn([e]) : J.push(e)), !i.is_fork && Dt.size > 0 && !kt && Nt();
	}
	return t;
}
function Nt() {
	kt = !1;
	for (let e of Dt) e.f & 1024 && A(e, h), jn(e) && In(e);
	Dt.clear();
}
function Pt(e) {
	F(e, e.v + 1);
}
function Ft(e, t, n) {
	var r = e.reactions;
	if (r !== null) for (var i = Pe(), a = r.length, o = 0; o < a; o++) {
		var s = r[o], c = s.f;
		if (!(!i && s === W)) {
			var l = (c & m) === 0;
			if (l && A(s, t), c & 2) {
				var u = s;
				M?.delete(u), c & 65536 || (c & 512 && (s.f |= C), Ft(u, h, n));
			} else if (l) {
				var d = s;
				c & 16 && rt !== null && rt.add(d), n === null ? st(d) : n.push(d);
			}
		}
	}
}
function It(t) {
	if (typeof t != "object" || !t || ie in t) return t;
	let n = c(t);
	if (n !== o && n !== s) return t;
	var r = /* @__PURE__ */ new Map(), i = e(t), l = /* @__PURE__ */ P(0), u = null, d = On, f = (e) => {
		if (On === d) return e();
		var t = H, n = On;
		U(null), kn(d);
		var r = e();
		return U(t), kn(n), r;
	};
	return i && r.set("length", /* @__PURE__ */ P(t.length, u)), new Proxy(t, {
		defineProperty(e, t, n) {
			(!("value" in n) || n.configurable === !1 || n.enumerable === !1 || n.writable === !1) && me();
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
					r.set(t, e), Pt(l);
				}
			} else F(n, w), Pt(l);
			return !0;
		},
		get(e, n, i) {
			if (n === ie) return t;
			var o = r.get(n), s = n in e;
			if (o === void 0 && (!s || a(e, n)?.writable) && (o = f(() => /* @__PURE__ */ P(It(s ? e[n] : w), u)), r.set(n, o)), o !== void 0) {
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
			if (t === ie) return !0;
			var n = r.get(t), i = n !== void 0 && n.v !== w || Reflect.has(e, t);
			return (n !== void 0 || W !== null && (!i || a(e, t)?.writable)) && (n === void 0 && (n = f(() => /* @__PURE__ */ P(i ? It(e[t]) : w, u)), r.set(t, n)), Y(n) === w) ? !1 : i;
		},
		set(e, t, n, o) {
			var s = r.get(t), c = t in e;
			if (i && t === "length") for (var d = n; d < s.v; d += 1) {
				var p = r.get(d + "");
				p === void 0 ? d in e && (p = f(() => /* @__PURE__ */ P(w, u)), r.set(d + "", p)) : F(p, w);
			}
			if (s === void 0) (!c || a(e, t)?.writable) && (s = f(() => /* @__PURE__ */ P(void 0, u)), F(s, It(n)), r.set(t, s));
			else {
				c = s.v !== w;
				var m = f(() => It(n));
				F(s, m);
			}
			var h = Reflect.getOwnPropertyDescriptor(e, t);
			if (h?.set && h.set.call(o, n), !c) {
				if (i && typeof t == "string") {
					var g = r.get("length"), _ = Number(t);
					Number.isInteger(_) && _ >= g.v && F(g, _ + 1);
				}
				Pt(l);
			}
			return !0;
		},
		ownKeys(e) {
			Y(l);
			var t = Reflect.ownKeys(e).filter((e) => {
				var t = r.get(e);
				return t === void 0 || t.v !== w;
			});
			for (var [n, i] of r) i.v !== w && !(n in e) && t.push(n);
			return t;
		},
		setPrototypeOf() {
			he();
		}
	});
}
var Lt, Rt, zt, Bt;
function Vt() {
	if (Lt === void 0) {
		Lt = window, document, Rt = /Firefox/.test(navigator.userAgent);
		var e = Element.prototype, t = Node.prototype, n = Text.prototype;
		zt = a(t, "firstChild").get, Bt = a(t, "nextSibling").get, l(e) && (e.__click = void 0, e.__className = void 0, e.__attributes = null, e.__style = void 0, e.__e = void 0), l(n) && (n.__t = void 0);
	}
}
function I(e = "") {
	return document.createTextNode(e);
}
/* @__NO_SIDE_EFFECTS__ */
function Ht(e) {
	return zt.call(e);
}
/* @__NO_SIDE_EFFECTS__ */
function Ut(e) {
	return Bt.call(e);
}
function L(e, t) {
	if (!T) return /* @__PURE__ */ Ht(e);
	var n = /* @__PURE__ */ Ht(E);
	if (n === null) n = E.appendChild(I());
	else if (t && n.nodeType !== 3) {
		var r = I();
		return n?.before(r), D(r), r;
	}
	return t && Jt(n), D(n), n;
}
function Wt(e, t = !1) {
	if (!T) {
		var n = /* @__PURE__ */ Ht(e);
		return n instanceof Comment && n.data === "" ? /* @__PURE__ */ Ut(n) : n;
	}
	if (t) {
		if (E?.nodeType !== 3) {
			var r = I();
			return E?.before(r), D(r), r;
		}
		Jt(E);
	}
	return E;
}
function R(e, t = 1, n = !1) {
	let r = T ? E : e;
	for (var i; t--;) i = r, r = /* @__PURE__ */ Ut(r);
	if (!T) return r;
	if (n) {
		if (r?.nodeType !== 3) {
			var a = I();
			return r === null ? i?.after(a) : r.before(a), D(a), a;
		}
		Jt(r);
	}
	return D(r), r;
}
function Gt(e) {
	e.textContent = "";
}
function Kt() {
	return !ke || rt !== null ? !1 : (W.f & v) !== 0;
}
function qt(e, t, n) {
	let r = n ? { is: n } : void 0;
	return document.createElementNS(t ?? "http://www.w3.org/1999/xhtml", e, r);
}
function Jt(e) {
	if (e.nodeValue.length < 65536) return;
	let t = e.nextSibling;
	for (; t !== null && t.nodeType === 3;) t.remove(), e.nodeValue += t.nodeValue, t = e.nextSibling;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/bindings/shared.js
function Yt(e) {
	var t = H, n = W;
	U(null), Cn(null);
	try {
		return e();
	} finally {
		U(t), Cn(n);
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/effects.js
function Xt(e) {
	W === null && (H === null && de(e), ue()), bn && le(e);
}
function Zt(e, t) {
	var n = t.last;
	n === null ? t.last = t.first = e : (n.next = e, e.prev = n, t.last = e);
}
function Qt(e, t) {
	var n = W;
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
	if (e & 4) Ze === null ? tt.ensure().schedule(r) : Ze.push(r);
	else if (t !== null) {
		try {
			In(r);
		} catch (e) {
			throw V(r), e;
		}
		i.deps === null && i.teardown === null && i.nodes === null && i.first === i.last && !(i.f & 524288) && (i = i.first, e & 16 && e & 65536 && i !== null && (i.f |= b));
	}
	if (i !== null && (i.parent = n, n !== null && Zt(i, n), H !== null && H.f & 2 && !(e & 64))) {
		var a = H;
		(a.effects ??= []).push(i);
	}
	return r;
}
function $t() {
	return H !== null && !Sn;
}
function en(e) {
	let t = Qt(8, null);
	return A(t, p), t.teardown = e, t;
}
function tn(e) {
	Xt("$effect");
	var t = W.f;
	if (!H && t & 32 && !(t & 32768)) {
		var n = k;
		(n.e ??= []).push(e);
	} else return nn(e);
}
function nn(e) {
	return Qt(4 | S, e);
}
function rn(e) {
	tt.ensure();
	let t = Qt(64 | x, e);
	return (e = {}) => new Promise((n) => {
		e.outro ? pn(t, () => {
			V(t), n(void 0);
		}) : (V(t), n(void 0));
	});
}
function an(e) {
	return Qt(ne | x, e);
}
function on(e, t = 0) {
	return Qt(8 | t, e);
}
function z(e, t = [], n = [], r = []) {
	mt(r, t, n, (t) => {
		Qt(8, () => e(...t.map(Y)));
	});
}
function sn(e, t = 0) {
	return Qt(16 | t, e);
}
function B(e) {
	return Qt(32 | x, e);
}
function cn(e) {
	var t = e.teardown;
	if (t !== null) {
		let e = bn, n = H;
		xn(!0), U(null);
		try {
			t.call(null);
		} finally {
			xn(e), U(n);
		}
	}
}
function ln(e, t = !1) {
	var n = e.first;
	for (e.first = e.last = null; n !== null;) {
		let e = n.ac;
		e !== null && Yt(() => {
			e.abort(oe);
		});
		var r = n.next;
		n.f & 64 ? n.parent = null : V(n, t), n = r;
	}
}
function un(e) {
	for (var t = e.first; t !== null;) {
		var n = t.next;
		t.f & 32 || V(t), t = n;
	}
}
function V(e, t = !0) {
	var n = !1;
	(t || e.f & 262144) && e.nodes !== null && e.nodes.end !== null && (dn(e.nodes.start, e.nodes.end), n = !0), A(e, y), ln(e, t && !n), Fn(e, 0);
	var r = e.nodes && e.nodes.t;
	if (r !== null) for (let e of r) e.stop();
	cn(e), e.f ^= y, e.f |= _;
	var i = e.parent;
	i !== null && i.first !== null && fn(e), e.next = e.prev = e.teardown = e.ctx = e.deps = e.fn = e.nodes = e.ac = null;
}
function dn(e, t) {
	for (; e !== null;) {
		var n = e === t ? null : /* @__PURE__ */ Ut(e);
		e.remove(), e = n;
	}
}
function fn(e) {
	var t = e.parent, n = e.prev, r = e.next;
	n !== null && (n.next = r), r !== null && (r.prev = n), t !== null && (t.first === e && (t.first = r), t.last === e && (t.last = n));
}
function pn(e, t, n = !0) {
	var r = [];
	mn(e, r, !0);
	var i = () => {
		n && V(e), t && t();
	}, a = r.length;
	if (a > 0) {
		var o = () => --a || i();
		for (var s of r) s.out(o);
	} else i();
}
function mn(e, t, n) {
	if (!(e.f & 8192)) {
		e.f ^= g;
		var r = e.nodes && e.nodes.t;
		if (r !== null) for (let e of r) (e.is_global || n) && t.push(e);
		for (var i = e.first; i !== null;) {
			var a = i.next, o = (i.f & 65536) != 0 || (i.f & 32) != 0 && (e.f & 16) != 0;
			mn(i, t, o ? n : !1), i = a;
		}
	}
}
function hn(e) {
	gn(e, !0);
}
function gn(e, t) {
	if (e.f & 8192) {
		e.f ^= g, e.f & 1024 || (A(e, m), tt.ensure().schedule(e));
		for (var n = e.first; n !== null;) {
			var r = n.next, i = (n.f & 65536) != 0 || (n.f & 32) != 0;
			gn(n, i ? t : !1), n = r;
		}
		var a = e.nodes && e.nodes.t;
		if (a !== null) for (let e of a) (e.is_global || t) && e.in();
	}
}
function _n(e, t) {
	if (e.nodes) for (var n = e.nodes.start, r = e.nodes.end; n !== null;) {
		var i = n === r ? null : /* @__PURE__ */ Ut(n);
		t.append(n), n = i;
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/legacy.js
var vn = null, yn = !1, bn = !1;
function xn(e) {
	bn = e;
}
var H = null, Sn = !1;
function U(e) {
	H = e;
}
var W = null;
function Cn(e) {
	W = e;
}
var G = null;
function wn(e) {
	H !== null && (!ke || H.f & 2) && (G === null ? G = [e] : G.push(e));
}
var K = null, q = 0, J = null;
function Tn(e) {
	J = e;
}
var En = 1, Dn = 0, On = Dn;
function kn(e) {
	On = e;
}
function An() {
	return ++En;
}
function jn(e) {
	var t = e.f;
	if (t & 2048) return !0;
	if (t & 2 && (e.f &= ~C), t & 4096) {
		for (var n = e.deps, r = n.length, i = 0; i < r; i++) {
			var a = n[i];
			if (jn(a) && wt(a), a.wv > e.wv) return !0;
		}
		t & 512 && M === null && A(e, p);
	}
	return !1;
}
function Mn(e, t, r = !0) {
	var i = e.reactions;
	if (i !== null && !(!ke && G !== null && n.call(G, e))) for (var a = 0; a < i.length; a++) {
		var o = i[a];
		o.f & 2 ? Mn(o, t, !1) : t === o && (r ? A(o, m) : o.f & 1024 && A(o, h), st(o));
	}
}
function Nn(e) {
	var t = K, n = q, r = J, i = H, a = G, o = k, s = Sn, c = On, l = e.f;
	K = null, q = 0, J = null, H = l & 96 ? null : e, G = null, je(e.ctx), Sn = !1, On = ++Dn, e.ac !== null && (Yt(() => {
		e.ac.abort(oe);
	}), e.ac = null);
	try {
		e.f |= te;
		var u = e.fn, d = u();
		e.f |= v;
		var f = e.deps, p = j?.is_fork;
		if (K !== null) {
			var m;
			if (p || Fn(e, q), f !== null && q > 0) for (f.length = q + K.length, m = 0; m < K.length; m++) f[q + m] = K[m];
			else e.deps = f = K;
			if ($t() && e.f & 512) for (m = q; m < f.length; m++) (f[m].reactions ??= []).push(e);
		} else !p && f !== null && q < f.length && (Fn(e, q), f.length = q);
		if (Pe() && J !== null && !Sn && f !== null && !(e.f & 6146)) for (m = 0; m < J.length; m++) Mn(J[m], e);
		if (i !== null && i !== e) {
			if (Dn++, i.deps !== null) for (let e = 0; e < n; e += 1) i.deps[e].rv = Dn;
			if (t !== null) for (let e of t) e.rv = Dn;
			J !== null && (r === null ? r = J : r.push(...J));
		}
		return e.f & 8388608 && (e.f ^= re), d;
	} catch (e) {
		return Re(e);
	} finally {
		e.f ^= te, K = t, q = n, J = r, H = i, G = a, je(o), Sn = s, On = c;
	}
}
function Pn(e, r) {
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
		s.f & 512 && (s.f ^= 512, s.f &= ~C), Ve(s), Tt(s), Fn(s, 0);
	}
}
function Fn(e, t) {
	var n = e.deps;
	if (n !== null) for (var r = t; r < n.length; r++) Pn(e, n[r]);
}
function In(e) {
	var t = e.f;
	if (!(t & 16384)) {
		A(e, p);
		var n = W, r = yn;
		W = e, yn = !0;
		try {
			t & 16777232 ? un(e) : ln(e), cn(e);
			var i = Nn(e);
			e.teardown = typeof i == "function" ? i : null, e.wv = En;
		} finally {
			yn = r, W = n;
		}
	}
}
function Y(e) {
	var t = (e.f & 2) != 0;
	if (vn?.add(e), H !== null && !Sn && !(W !== null && W.f & 16384) && (G === null || !n.call(G, e))) {
		var r = H.deps;
		if (H.f & 2097152) e.rv < Dn && (e.rv = Dn, K === null && r !== null && r[q] === e ? q++ : K === null ? K = [e] : K.push(e));
		else {
			(H.deps ??= []).push(e);
			var i = e.reactions;
			i === null ? e.reactions = [H] : n.call(i, H) || i.push(H);
		}
	}
	if (bn && Ot.has(e)) return Ot.get(e);
	if (t) {
		var a = e;
		if (bn) {
			var o = a.v;
			return (!(a.f & 1024) && a.reactions !== null || Rn(a)) && (o = Ct(a)), Ot.set(a, o), o;
		}
		var s = (a.f & 512) == 0 && !Sn && H !== null && (yn || (H.f & 512) != 0), c = (a.f & v) === 0;
		jn(a) && (s && (a.f |= 512), wt(a)), s && !c && (Et(a), Ln(a));
	}
	if (M?.has(e)) return M.get(e);
	if (e.f & 8388608) throw e.v;
	return e.v;
}
function Ln(e) {
	if (e.f |= 512, e.deps !== null) for (let t of e.deps) (t.reactions ??= []).push(e), t.f & 2 && !(t.f & 512) && (Et(t), Ln(t));
}
function Rn(e) {
	if (e.v === w) return !0;
	if (e.deps === null) return !1;
	for (let t of e.deps) if (Ot.has(t) || t.f & 2 && Rn(t)) return !0;
	return !1;
}
function zn(e) {
	var t = Sn;
	try {
		return Sn = !0, e();
	} finally {
		Sn = t;
	}
}
[.../* @__PURE__ */ "allowfullscreen.async.autofocus.autoplay.checked.controls.default.disabled.formnovalidate.indeterminate.inert.ismap.loop.multiple.muted.nomodule.novalidate.open.playsinline.readonly.required.reversed.seamless.selected.webkitdirectory.defer.disablepictureinpicture.disableremoteplayback".split(".")];
var Bn = ["touchstart", "touchmove"];
function Vn(e) {
	return Bn.includes(e);
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/events.js
var Hn = Symbol("events"), Un = /* @__PURE__ */ new Set(), Wn = /* @__PURE__ */ new Set();
function Gn(e, t, n) {
	(t[Hn] ??= {})[e] = n;
}
function Kn(e) {
	for (var t = 0; t < e.length; t++) Un.add(e[t]);
	for (var n of Wn) n(e);
}
var qn = null;
function Jn(e) {
	var t = this, n = t.ownerDocument, r = e.type, a = e.composedPath?.() || [], o = a[0] || e.target;
	qn = e;
	var s = 0, c = qn === e && e[Hn];
	if (c) {
		var l = a.indexOf(c);
		if (l !== -1 && (t === document || t === window)) {
			e[Hn] = t;
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
		U(null), Cn(null);
		try {
			for (var p, m = []; o !== null;) {
				var h = o.assignedSlot || o.parentNode || o.host || null;
				try {
					var g = o[Hn]?.[r];
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
			e[Hn] = t, delete e.currentTarget, U(d), Cn(f);
		}
	}
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/reconciler.js
var Yn = globalThis?.window?.trustedTypes && /* @__PURE__ */ globalThis.window.trustedTypes.createPolicy("svelte-trusted-html", { createHTML: (e) => e });
function Xn(e) {
	return Yn?.createHTML(e) ?? e;
}
function Zn(e) {
	var t = qt("template");
	return t.innerHTML = Xn(e.replaceAll("<!>", "<!---->")), t.content;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/template.js
function Qn(e, t) {
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
		if (T) return Qn(E, null), E;
		i === void 0 && (i = Zn(a ? e : "<!>" + e), n || (i = /* @__PURE__ */ Ht(i)));
		var t = r || Rt ? document.importNode(i, !0) : i.cloneNode(!0);
		if (n) {
			var o = /* @__PURE__ */ Ht(t), s = t.lastChild;
			Qn(o, s);
		} else Qn(t, t);
		return t;
	};
}
function $n() {
	if (T) return Qn(E, null), E;
	var e = document.createDocumentFragment(), t = document.createComment(""), n = I();
	return e.append(t, n), Qn(t, n), e;
}
function Z(e, t) {
	if (T) {
		var n = W;
		(!(n.f & 32768) || n.nodes.end === null) && (n.nodes.end = E), Se();
		return;
	}
	e !== null && e.before(t);
}
function Q(e, t) {
	var n = t == null ? "" : typeof t == "object" ? `${t}` : t;
	n !== (e.__t ??= e.nodeValue) && (e.__t = n, e.nodeValue = `${n}`);
}
function er(e, t) {
	return nr(e, t);
}
var tr = /* @__PURE__ */ new Map();
function nr(e, { target: t, anchor: n, props: i = {}, events: a, context: o, intro: s = !0, transformError: c }) {
	Vt();
	var l = void 0, u = rn(() => {
		var s = n ?? t.appendChild(I());
		ft(s, { pending: () => {} }, (t) => {
			Me({});
			var n = k;
			if (o && (n.c = o), a && (i.$$events = a), T && Qn(t, null), l = e(t, i) || {}, T && (W.nodes.end = E, E === null || E.nodeType !== 8 || E.data !== "]")) throw ye(), ve;
			Ne();
		}, c);
		var u = /* @__PURE__ */ new Set(), d = (e) => {
			for (var n = 0; n < e.length; n++) {
				var r = e[n];
				if (!u.has(r)) {
					u.add(r);
					var i = Vn(r);
					for (let e of [t, document]) {
						var a = tr.get(e);
						a === void 0 && (a = /* @__PURE__ */ new Map(), tr.set(e, a));
						var o = a.get(r);
						o === void 0 ? (e.addEventListener(r, Jn, { passive: i }), a.set(r, 1)) : a.set(r, o + 1);
					}
				}
			}
		};
		return d(r(Un)), Wn.add(d), () => {
			for (var e of u) for (let n of [t, document]) {
				var r = tr.get(n), i = r.get(e);
				--i == 0 ? (n.removeEventListener(e, Jn), r.delete(e), r.size === 0 && tr.delete(n)) : r.set(e, i);
			}
			Wn.delete(d), s !== n && s.parentNode?.removeChild(s);
		};
	});
	return rr.set(l, u), l;
}
var rr = /* @__PURE__ */ new WeakMap(), ir = class {
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
			if (n) hn(n), this.#r.delete(t);
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
						_n(r, t), t.append(I()), this.#n.set(e, {
							effect: r,
							fragment: t
						});
					} else V(r);
					this.#r.delete(e), this.#t.delete(e);
				};
				this.#i || !n ? (this.#r.add(e), pn(r, i, !1)) : i();
			}
		}
	};
	#o = (e) => {
		this.#e.delete(e);
		let t = Array.from(this.#e.values());
		for (let [e, n] of this.#n) t.includes(e) || (V(n.effect), this.#n.delete(e));
	};
	ensure(e, t) {
		var n = j, r = Kt();
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
	T && (r = E, Se());
	var i = new ir(e), a = n ? b : 0;
	function o(e, t) {
		if (T) {
			var n = Te(r);
			if (e !== parseInt(n.substring(1))) {
				var a = we();
				D(a), i.anchor = a, xe(!1), i.ensure(e, t), xe(!0);
				return;
			}
		}
		i.ensure(e, t);
	}
	sn(() => {
		var e = !1;
		t((t, n = 0) => {
			e = !0, o(n, t);
		}), e || o(-1, null);
	}, a);
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/blocks/each.js
function ar(e, t) {
	return t;
}
function or(e, t, n) {
	for (var i = [], a = t.length, o, s = t.length, c = 0; c < a; c++) {
		let n = t[c];
		pn(n, () => {
			if (o) {
				if (o.pending.delete(n), o.done.add(n), o.pending.size === 0) {
					var t = e.outrogroups;
					sr(e, r(o.done)), t.delete(o), t.size === 0 && (e.outrogroups = null);
				}
			} else --s;
		}, !1);
	}
	if (s === 0) {
		var l = i.length === 0 && n !== null;
		if (l) {
			var u = n, d = u.parentNode;
			Gt(d), d.append(u), e.items.clear();
		}
		sr(e, t, !l);
	} else o = {
		pending: new Set(t),
		done: /* @__PURE__ */ new Set()
	}, (e.outrogroups ??= /* @__PURE__ */ new Set()).add(o);
}
function sr(e, t, n = !0) {
	var r;
	if (e.pending.size > 0) {
		r = /* @__PURE__ */ new Set();
		for (let t of e.pending.values()) for (let n of t) r.add(e.items.get(n).e);
	}
	for (var i = 0; i < t.length; i++) {
		var a = t[i];
		r?.has(a) ? (a.f |= ee, _n(a, document.createDocumentFragment())) : V(t[i], n);
	}
}
var cr;
function lr(t, n, i, a, o, s = null) {
	var c = t, l = /* @__PURE__ */ new Map();
	if (n & 4) {
		var u = t;
		c = T ? D(/* @__PURE__ */ Ht(u)) : u.appendChild(I());
	}
	T && Se();
	var d = null, f = /* @__PURE__ */ bt(() => {
		var t = i();
		return e(t) ? t : t == null ? [] : r(t);
	}), p, m = /* @__PURE__ */ new Map(), h = !0;
	function g(e) {
		v.effect.f & 16384 || (v.pending.delete(e), v.fallback = d, dr(v, p, c, n, a), d !== null && (p.length === 0 ? d.f & 33554432 ? (d.f ^= ee, pr(d, null, c)) : hn(d) : pn(d, () => {
			d = null;
		})));
	}
	function _(e) {
		v.pending.delete(e);
	}
	var v = {
		effect: sn(() => {
			p = Y(f);
			var e = p.length;
			let t = !1;
			T && Te(c) === "[!" != (e === 0) && (c = we(), D(c), xe(!1), t = !0);
			for (var r = /* @__PURE__ */ new Set(), u = j, v = Kt(), y = 0; y < e; y += 1) {
				T && E.nodeType === 8 && E.data === "]" && (c = E, t = !0, xe(!1));
				var b = p[y], x = a(b, y), S = h ? null : l.get(x);
				S ? (S.v && Mt(S.v, b), S.i && Mt(S.i, y), v && u.unskip_effect(S.e)) : (S = fr(l, h ? c : cr ??= I(), b, x, y, o, n, i), h || (S.e.f |= ee), l.set(x, S)), r.add(x);
			}
			if (e === 0 && s && !d && (h ? d = B(() => s(c)) : (d = B(() => s(cr ??= I())), d.f |= ee)), e > r.size && ce("", "", ""), T && e > 0 && D(we()), !h) if (m.set(u, r), v) {
				for (let [e, t] of l) r.has(e) || u.skip_effect(t.e);
				u.oncommit(g), u.ondiscard(_);
			} else g(u);
			t && xe(!0), Y(f);
		}),
		flags: n,
		items: l,
		pending: m,
		outrogroups: null,
		fallback: d
	};
	h = !1, T && (c = E);
}
function ur(e) {
	for (; e !== null && !(e.f & 32);) e = e.next;
	return e;
}
function dr(e, t, n, i, a) {
	var o = (i & 8) != 0, s = t.length, c = e.items, l = ur(e.effect.first), u, d = null, f, p = [], m = [], h, g, _, v;
	if (o) for (v = 0; v < s; v += 1) h = t[v], g = a(h, v), _ = c.get(g).e, _.f & 33554432 || (_.nodes?.a?.measure(), (f ??= /* @__PURE__ */ new Set()).add(_));
	for (v = 0; v < s; v += 1) {
		if (h = t[v], g = a(h, v), _ = c.get(g).e, e.outrogroups !== null) for (let t of e.outrogroups) t.pending.delete(_), t.done.delete(_);
		if (_.f & 33554432) if (_.f ^= ee, _ === l) pr(_, null, n);
		else {
			var y = d ? d.next : l;
			_ === e.effect.last && (e.effect.last = _.prev), _.prev && (_.prev.next = _.next), _.next && (_.next.prev = _.prev), mr(e, d, _), mr(e, _, y), pr(_, y, n), d = _, p = [], m = [], l = ur(d.next);
			continue;
		}
		if (_.f & 8192 && (hn(_), o && (_.nodes?.a?.unfix(), (f ??= /* @__PURE__ */ new Set()).delete(_))), _ !== l) {
			if (u !== void 0 && u.has(_)) {
				if (p.length < m.length) {
					var b = m[0], x;
					d = b.prev;
					var S = p[0], C = p[p.length - 1];
					for (x = 0; x < p.length; x += 1) pr(p[x], b, n);
					for (x = 0; x < m.length; x += 1) u.delete(m[x]);
					mr(e, S.prev, C.next), mr(e, d, S), mr(e, C, b), l = b, d = C, --v, p = [], m = [];
				} else u.delete(_), pr(_, l, n), mr(e, _.prev, _.next), mr(e, _, d === null ? e.effect.first : d.next), mr(e, d, _), d = _;
				continue;
			}
			for (p = [], m = []; l !== null && l !== _;) (u ??= /* @__PURE__ */ new Set()).add(l), m.push(l), l = ur(l.next);
			if (l === null) continue;
		}
		_.f & 33554432 || p.push(_), d = _, l = ur(_.next);
	}
	if (e.outrogroups !== null) {
		for (let t of e.outrogroups) t.pending.size === 0 && (sr(e, r(t.done)), e.outrogroups?.delete(t));
		e.outrogroups.size === 0 && (e.outrogroups = null);
	}
	if (l !== null || u !== void 0) {
		var te = [];
		if (u !== void 0) for (_ of u) _.f & 8192 || te.push(_);
		for (; l !== null;) !(l.f & 8192) && l !== e.fallback && te.push(l), l = ur(l.next);
		var ne = te.length;
		if (ne > 0) {
			var re = i & 4 && s === 0 ? n : null;
			if (o) {
				for (v = 0; v < ne; v += 1) te[v].nodes?.a?.measure();
				for (v = 0; v < ne; v += 1) te[v].nodes?.a?.fix();
			}
			or(e, te, re);
		}
	}
	o && Le(() => {
		if (f !== void 0) for (_ of f) _.nodes?.a?.apply();
	});
}
function fr(e, t, n, r, i, a, o, s) {
	var c = o & 1 ? o & 16 ? At(n) : /* @__PURE__ */ jt(n, !1, !1) : null, l = o & 2 ? At(i) : null;
	return {
		v: c,
		i: l,
		e: B(() => (a(t, c ?? n, l ?? i, s), () => {
			e.delete(r);
		}))
	};
}
function pr(e, t, n) {
	if (e.nodes) for (var r = e.nodes.start, i = e.nodes.end, a = t && !(t.f & 33554432) ? t.nodes.start : n; r !== null;) {
		var o = /* @__PURE__ */ Ut(r);
		if (a.before(r), r === i) return;
		r = o;
	}
}
function mr(e, t, n) {
	t === null ? e.effect.first = n : t.next = n, n === null ? e.effect.last = t : n.prev = t;
}
//#endregion
//#region node_modules/clsx/dist/clsx.mjs
function hr(e) {
	var t, n, r = "";
	if (typeof e == "string" || typeof e == "number") r += e;
	else if (typeof e == "object") if (Array.isArray(e)) {
		var i = e.length;
		for (t = 0; t < i; t++) e[t] && (n = hr(e[t])) && (r && (r += " "), r += n);
	} else for (n in e) e[n] && (r && (r += " "), r += n);
	return r;
}
function gr() {
	for (var e, t, n = 0, r = "", i = arguments.length; n < i; n++) (e = arguments[n]) && (t = hr(e)) && (r && (r += " "), r += t);
	return r;
}
//#endregion
//#region node_modules/svelte/src/internal/shared/attributes.js
function _r(e) {
	return typeof e == "object" ? gr(e) : e ?? "";
}
var vr = [..." 	\n\r\f\xA0\v﻿"];
function yr(e, t, n) {
	var r = e == null ? "" : "" + e;
	if (t && (r = r ? r + " " + t : t), n) {
		for (var i of Object.keys(n)) if (n[i]) r = r ? r + " " + i : i;
		else if (r.length) for (var a = i.length, o = 0; (o = r.indexOf(i, o)) >= 0;) {
			var s = o + a;
			(o === 0 || vr.includes(r[o - 1])) && (s === r.length || vr.includes(r[s])) ? r = (o === 0 ? "" : r.substring(0, o)) + r.substring(s + 1) : o = s;
		}
	}
	return r === "" ? null : r;
}
//#endregion
//#region node_modules/svelte/src/internal/client/dom/elements/class.js
function br(e, t, n, r, i, a) {
	var o = e.__className;
	if (T || o !== n || o === void 0) {
		var s = yr(n, r, a);
		(!T || s !== e.getAttribute("class")) && (s == null ? e.removeAttribute("class") : t ? e.className = s : e.setAttribute("class", s)), e.__className = n;
	} else if (a && i !== a) for (var c in a) {
		var l = !!a[c];
		(i == null || l !== !!i[c]) && e.classList.toggle(c, l);
	}
	return a;
}
//#endregion
//#region node_modules/svelte/src/internal/client/reactivity/props.js
function xr(e, t, n, r) {
	var i = !Ae || (n & 2) != 0, o = (n & 8) != 0, s = (n & 16) != 0, c = r, l = !0, u = () => (l && (l = !1, c = s ? zn(r) : r), c);
	let d;
	if (o) {
		var f = ie in e || ae in e;
		d = a(e, t)?.set ?? (f && t in e ? (n) => e[t] = n : void 0);
	}
	var p, m = !1;
	o ? [p, m] = Ke(() => e[t]) : p = e[t], p === void 0 && r !== void 0 && (p = u(), d && (i && pe(t), d(p)));
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
	var _ = !1, v = (n & 1 ? vt : bt)(() => (_ = !1, h()));
	o && Y(v);
	var y = W;
	return (function(e, t) {
		if (arguments.length > 0) {
			let n = t ? Y(v) : i && o ? It(e) : e;
			return F(v, n), _ = !0, c !== void 0 && (c = n), e;
		}
		return bn && _ || y.f & 16384 ? v.v : Y(v);
	});
}
//#endregion
//#region node_modules/svelte/src/internal/disclose-version.js
typeof window < "u" && ((window.__svelte ??= {}).v ??= /* @__PURE__ */ new Set()).add("5");
//#endregion
//#region src/lib/markup.js
function Sr(e, t, n, r) {
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
function Cr(e, t) {
	return e * (1 + t / 100);
}
function wr(e, t, n) {
	return e * Cr(t, n);
}
function Tr(e, t) {
	let n = {
		materials: 0,
		labor: 0,
		equipment: 0,
		subs: 0,
		other: 0
	}, r = 0, i = 0, a = (a) => {
		let o = Sr(a.category_type, t, e.markup_overrides, e.markup_enabled), s = a.quantity * a.unit_price, c = wr(a.quantity, a.unit_price, o);
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
function Er(e, t) {
	let n = {
		materials: 0,
		labor: 0,
		equipment: 0,
		subs: 0,
		other: 0
	}, r = 0, i = 0;
	for (let a of e.subcategories) {
		let e = Tr(a, t);
		r += e.base, i += e.withMarkup;
		for (let t of Object.keys(n)) n[t] += e.byType[t];
	}
	return {
		base: r,
		withMarkup: i,
		byType: n
	};
}
function Dr(e) {
	let t = {
		materials: 0,
		labor: 0,
		equipment: 0,
		subs: 0,
		other: 0
	}, n = 0, r = 0;
	for (let i of e.sections) {
		let a = Er(i, e.globals);
		n += a.base, r += a.withMarkup;
		for (let e of Object.keys(t)) t[e] += a.byType[e];
	}
	return {
		base: n,
		withMarkup: r,
		byType: t
	};
}
function Or(e) {
	return new Intl.NumberFormat("en-US", {
		style: "currency",
		currency: "USD",
		minimumFractionDigits: 2
	}).format(e);
}
function kr(e) {
	return `${e}%`;
}
//#endregion
//#region src/lib/LineItemRow.svelte
var Ar = /* @__PURE__ */ X("<span class=\"block text-xs text-slate-400 mt-0.5\"> </span>"), jr = /* @__PURE__ */ X("<tr class=\"border-b border-slate-100 hover:bg-slate-50 text-sm\"><td class=\"px-2 py-1.5 w-20\"><span> </span></td><td class=\"px-2 py-1.5\"><span class=\"text-slate-800\"> </span> <!></td><td class=\"px-2 py-1.5 text-right font-mono w-20\"> </td><td class=\"px-2 py-1.5 text-center text-slate-500 w-16\"> </td><td class=\"px-2 py-1.5 text-right font-mono w-24\"> </td><td class=\"px-2 py-1.5 text-right font-mono text-slate-500 w-16\"> </td><td class=\"px-2 py-1.5 text-right font-mono w-24\"> </td><td class=\"px-2 py-1.5 text-right font-mono font-medium w-28\"> </td></tr>");
function Mr(e, t) {
	Me(t, !0);
	let n = {
		materials: "bg-blue-100 text-blue-700",
		labor: "bg-amber-100 text-amber-700",
		equipment: "bg-purple-100 text-purple-700",
		subs: "bg-green-100 text-green-700",
		other: "bg-slate-100 text-slate-600"
	}, r = /* @__PURE__ */ N(() => Sr(t.item.category_type, t.globals, t.markupOverrides, t.markupEnabled)), i = /* @__PURE__ */ N(() => Cr(t.item.unit_price, Y(r))), a = /* @__PURE__ */ N(() => wr(t.item.quantity, t.item.unit_price, Y(r))), o = /* @__PURE__ */ N(() => n[t.item.category_type] || n.other);
	var s = jr(), c = L(s), l = L(c), u = L(l, !0);
	O(l), O(c);
	var d = R(c), f = L(d), p = L(f, !0);
	O(f);
	var m = R(f, 2), h = (e) => {
		var n = Ar(), r = L(n, !0);
		O(n), z(() => Q(r, t.item.description)), Z(e, n);
	};
	$(m, (e) => {
		t.item.description && e(h);
	}), O(d);
	var g = R(d), _ = L(g, !0);
	O(g);
	var v = R(g), y = L(v, !0);
	O(v);
	var b = R(v), x = L(b, !0);
	O(b);
	var S = R(b), ee = L(S);
	O(S);
	var C = R(S), te = L(C, !0);
	O(C);
	var ne = R(C), re = L(ne, !0);
	O(ne), O(s), z((e, n, i) => {
		br(l, 1, `text-xs px-1.5 py-0.5 rounded font-medium ${Y(o) ?? ""}`), Q(u, t.item.category_type), Q(p, t.item.item_name), Q(_, t.item.quantity), Q(y, t.item.unit), Q(x, e), Q(ee, `${Y(r) ?? ""}%`), Q(te, n), Q(re, i);
	}, [
		() => Or(t.item.unit_price),
		() => Or(Y(i)),
		() => Or(Y(a))
	]), Z(e, s), Ne();
}
//#endregion
//#region src/lib/ComponentGroupBlock.svelte
var Nr = /* @__PURE__ */ X("<table class=\"w-full\"><tbody></tbody></table>"), Pr = /* @__PURE__ */ X("<div class=\"ml-4 mt-2\"><div class=\"flex items-center gap-2 mb-1\"><span class=\"text-xs font-semibold text-slate-500 uppercase tracking-wide\"> </span> <span class=\"text-xs text-slate-400\"> </span></div> <!></div>");
function Fr(e, t) {
	Me(t, !0);
	var n = Pr(), r = L(n), i = L(r), a = L(i, !0);
	O(i);
	var o = R(i, 2), s = L(o);
	O(o), O(r);
	var c = R(r, 2), l = (e) => {
		var n = Nr(), r = L(n);
		lr(r, 21, () => t.group.line_items, (e) => e.id, (e, n) => {
			Mr(e, {
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
				}
			});
		}), O(r), O(n), Z(e, n);
	};
	$(c, (e) => {
		t.group.line_items.length > 0 && e(l);
	}), O(n), z(() => {
		Q(a, t.group.name), Q(s, `(${t.group.line_items.length ?? ""})`);
	}), Z(e, n), Ne();
}
//#endregion
//#region src/lib/SubcategoryBlock.svelte
var Ir = /* @__PURE__ */ X("<span class=\"text-xs px-1.5 py-0.5 rounded bg-amber-50 text-amber-600 font-medium\">overrides</span>"), Lr = /* @__PURE__ */ X("<span class=\"text-xs px-1.5 py-0.5 rounded bg-green-50 text-green-600 font-medium\"> </span>"), Rr = /* @__PURE__ */ X("<span> </span>"), zr = /* @__PURE__ */ X("<div class=\"flex items-center gap-3 py-1.5 px-2 bg-amber-50 rounded text-xs mb-2\"><span class=\"font-medium text-amber-700\">Markup:</span> <!></div>"), Br = /* @__PURE__ */ X("<table class=\"w-full\"><thead><tr class=\"text-xs text-slate-400 uppercase tracking-wide border-b border-slate-100\"><th class=\"px-2 py-1 text-left w-20\">Type</th><th class=\"px-2 py-1 text-left\">Name</th><th class=\"px-2 py-1 text-right w-20\">Qty</th><th class=\"px-2 py-1 text-center w-16\">Unit</th><th class=\"px-2 py-1 text-right w-24\">Price</th><th class=\"px-2 py-1 text-right w-16\">Markup</th><th class=\"px-2 py-1 text-right w-24\">w/ Markup</th><th class=\"px-2 py-1 text-right w-28\">Total</th></tr></thead><tbody></tbody></table>"), Vr = /* @__PURE__ */ X("<span class=\"text-green-600\"> </span>"), Hr = /* @__PURE__ */ X("<div class=\"px-4 pb-3\"><!> <!> <!> <div class=\"flex justify-between items-center mt-2 pt-2 border-t border-slate-100 text-sm\"><span class=\"text-slate-500\"> <!></span> <span class=\"font-mono font-semibold text-slate-700\"> </span></div></div>"), Ur = /* @__PURE__ */ X("<div class=\"border-t border-slate-200\"><button class=\"w-full flex items-center justify-between px-4 py-2 hover:bg-slate-50 transition-colors text-left\"><div class=\"flex items-center gap-2\"><svg fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> <span class=\"font-medium text-slate-600 text-sm\"> </span> <span class=\"text-xs text-slate-400\"> </span> <!> <!></div> <span class=\"font-mono text-sm font-semibold text-slate-700\"> </span></button> <!></div>");
function Wr(e, t) {
	Me(t, !0);
	let n = /* @__PURE__ */ P(It(xr(t, "collapsed", 3, !1)())), r = /* @__PURE__ */ N(() => Tr(t.subcat, t.globals)), i = /* @__PURE__ */ N(() => [
		"materials",
		"labor",
		"equipment",
		"subs",
		"other"
	].map((e) => ({
		type: e,
		value: Sr(e, t.globals, t.subcat.markup_overrides, t.subcat.markup_enabled),
		isOverride: t.subcat.markup_overrides[e] != null,
		isDisabled: !t.subcat.markup_enabled[e]
	}))), a = /* @__PURE__ */ N(() => Y(i).some((e) => e.isOverride || e.isDisabled)), o = /* @__PURE__ */ N(() => t.subcat.line_items.length + t.subcat.component_groups.reduce((e, t) => e + t.line_items.length, 0));
	var s = Ur(), c = L(s), l = L(c), u = L(l), d = R(u, 2), f = L(d, !0);
	O(d);
	var p = R(d, 2), m = L(p);
	O(p);
	var h = R(p, 2), g = (e) => {
		Z(e, Ir());
	};
	$(h, (e) => {
		Y(a) && e(g);
	});
	var _ = R(h, 2), v = (e) => {
		var n = Lr(), r = L(n);
		O(n), z((e) => Q(r, `+${e ?? ""} lump sum`), [() => Or(t.subcat.lump_sum)]), Z(e, n);
	};
	$(_, (e) => {
		t.subcat.lump_sum > 0 && e(v);
	}), O(l);
	var y = R(l, 2), b = L(y, !0);
	O(y), O(c);
	var x = R(c, 2), S = (e) => {
		var n = Hr(), s = L(n), c = (e) => {
			var t = zr();
			lr(R(L(t), 2), 17, () => Y(i), ar, (e, t) => {
				var n = Rr(), r = L(n);
				O(n), z((e) => {
					br(n, 1, _r(Y(t).isDisabled ? "text-slate-400 line-through" : Y(t).isOverride ? "text-amber-700 font-medium" : "text-slate-500")), Q(r, `${Y(t).type ?? ""} ${e ?? ""}`);
				}, [() => kr(Y(t).value)]), Z(e, n);
			}), O(t), Z(e, t);
		};
		$(s, (e) => {
			Y(a) && e(c);
		});
		var l = R(s, 2), u = (e) => {
			var n = Br(), r = R(L(n));
			lr(r, 21, () => t.subcat.line_items, (e) => e.id, (e, n) => {
				Mr(e, {
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
					}
				});
			}), O(r), O(n), Z(e, n);
		};
		$(l, (e) => {
			Y(o) > 0 && e(u);
		});
		var d = R(l, 2);
		lr(d, 17, () => t.subcat.component_groups, (e) => e.id, (e, n) => {
			Fr(e, {
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
				}
			});
		});
		var f = R(d, 2), p = L(f), m = L(p), h = R(m), g = (e) => {
			var n = Vr(), r = L(n);
			O(n), z((e) => Q(r, `+ ${e ?? ""} lump sum`), [() => Or(t.subcat.lump_sum)]), Z(e, n);
		};
		$(h, (e) => {
			t.subcat.lump_sum > 0 && e(g);
		}), O(p);
		var _ = R(p, 2), v = L(_, !0);
		O(_), O(f), O(n), z((e, t) => {
			Q(m, `Subtotal: ${e ?? ""} `), Q(v, t);
		}, [() => Or(Y(r).base), () => Or(Y(r).withMarkup)]), Z(e, n);
	};
	$(x, (e) => {
		Y(n) || e(S);
	}), O(s), z((e) => {
		br(u, 0, `w-4 h-4 text-slate-400 transition-transform ${Y(n) ? "" : "rotate-90"}`), Q(f, t.subcat.name), Q(m, `${Y(o) ?? ""} item${Y(o) === 1 ? "" : "s"}`), Q(b, e);
	}, [() => Or(Y(r).withMarkup)]), Gn("click", c, () => F(n, !Y(n))), Z(e, s), Ne();
}
Kn(["click"]);
//#endregion
//#region src/lib/SectionBlock.svelte
var Gr = /* @__PURE__ */ X("<div class=\"px-4 py-8 text-center text-slate-400 text-sm\">No subcategories yet</div>"), Kr = /* @__PURE__ */ X("<div class=\"mb-4 border border-slate-200 rounded-lg overflow-hidden bg-white\"><button class=\"w-full flex items-center justify-between px-4 py-3 bg-slate-800 text-white hover:bg-slate-700 transition-colors text-left\"><div class=\"flex items-center gap-3\"><svg fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 5l7 7-7 7\"></path></svg> <span class=\"font-semibold\"> </span> <span class=\"text-xs text-slate-400\"> </span></div> <span class=\"font-mono font-semibold\"> </span></button> <!></div>");
function qr(e, t) {
	Me(t, !0);
	let n = /* @__PURE__ */ P(It(xr(t, "collapsed", 3, !1)())), r = /* @__PURE__ */ N(() => Er(t.section, t.globals)), i = /* @__PURE__ */ N(() => t.section.subcategories.reduce((e, t) => e + t.line_items.length + t.component_groups.reduce((e, t) => e + t.line_items.length, 0), 0));
	var a = Kr(), o = L(a), s = L(o), c = L(s), l = R(c, 2), u = L(l, !0);
	O(l);
	var d = R(l, 2), f = L(d);
	O(d), O(s);
	var p = R(s, 2), m = L(p, !0);
	O(p), O(o);
	var h = R(o, 2), g = (e) => {
		var n = $n(), r = Wt(n), i = (e) => {
			Z(e, Gr());
		}, a = (e) => {
			var n = $n();
			lr(Wt(n), 17, () => t.section.subcategories, (e) => e.id, (e, n) => {
				Wr(e, {
					get subcat() {
						return Y(n);
					},
					get globals() {
						return t.globals;
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
		br(c, 0, `w-4 h-4 text-slate-400 transition-transform ${Y(n) ? "" : "rotate-90"}`), Q(u, t.section.name), Q(f, `${t.section.subcategories.length ?? ""} subcategor${t.section.subcategories.length === 1 ? "y" : "ies"}
				· ${Y(i) ?? ""} item${Y(i) === 1 ? "" : "s"}`), Q(m, e);
	}, [() => Or(Y(r).withMarkup)]), Gn("click", o, () => F(n, !Y(n))), Z(e, a), Ne();
}
Kn(["click"]);
//#endregion
//#region src/lib/FooterSummary.svelte
var Jr = /* @__PURE__ */ X("<div class=\"flex flex-col\"><span class=\"text-xs text-slate-400 uppercase tracking-wide\"> </span> <span class=\"font-mono\"> </span></div>"), Yr = /* @__PURE__ */ X("<span class=\"text-slate-500 text-xs\">No items yet</span>"), Xr = /* @__PURE__ */ X("<div class=\"fixed bottom-0 left-0 right-0 bg-slate-900 text-white border-t border-slate-700 z-20\"><div class=\"flex items-center justify-between px-4 py-3\"><div class=\"flex items-center gap-6 text-sm\"><div class=\"flex flex-col\"><span class=\"text-xs text-slate-400 uppercase tracking-wide\">Base Cost</span> <span class=\"font-mono\"> </span></div> <div class=\"w-px h-8 bg-slate-700\"></div> <!> <!></div> <div class=\"flex flex-col items-end\"><span class=\"text-xs text-slate-400 uppercase tracking-wide\">Total</span> <span class=\"font-mono text-lg font-bold\"> </span></div></div></div>");
function Zr(e, t) {
	Me(t, !0);
	let n = /* @__PURE__ */ N(() => Dr(t.estimate)), r = /* @__PURE__ */ N(() => Object.entries(Y(n).byType).filter(([, e]) => e > 0).map(([e, t]) => ({
		type: e,
		value: t
	}))), i = {
		materials: "Materials",
		labor: "Labor",
		equipment: "Equipment",
		subs: "Subs",
		other: "Other"
	};
	var a = Xr(), o = L(a), s = L(o), c = L(s), l = R(L(c), 2), u = L(l, !0);
	O(l), O(c);
	var d = R(c, 4);
	lr(d, 17, () => Y(r), ar, (e, t) => {
		let n = () => Y(t).type, r = () => Y(t).value;
		var a = Jr(), o = L(a), s = L(o, !0);
		O(o);
		var c = R(o, 2), l = L(c, !0);
		O(c), O(a), z((e) => {
			Q(s, i[n()]), Q(l, e);
		}, [() => Or(r())]), Z(e, a);
	});
	var f = R(d, 2), p = (e) => {
		Z(e, Yr());
	};
	$(f, (e) => {
		Y(r).length === 0 && e(p);
	}), O(s);
	var m = R(s, 2), h = R(L(m), 2), g = L(h, !0);
	O(h), O(m), O(o), O(a), z((e, t) => {
		Q(u, e), Q(g, t);
	}, [() => Or(Y(n).base), () => Or(Y(n).withMarkup)]), Z(e, a), Ne();
}
//#endregion
//#region src/EstimateBuilder.svelte
var Qr = /* @__PURE__ */ X("<div class=\"flex items-center justify-center h-64\"><div class=\"text-slate-500\">Loading estimate...</div></div>"), $r = /* @__PURE__ */ X("<div class=\"flex items-center justify-center h-64\"><div class=\"text-red-500\"> </div></div>"), ei = /* @__PURE__ */ X("<div class=\"text-center py-16 text-slate-400\"><svg class=\"w-12 h-12 mx-auto mb-3 text-slate-300\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"1.5\" d=\"M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z\"></path></svg> <p class=\"text-lg font-medium\">No sections yet</p> <p class=\"text-sm mt-1\">Add a section to start building your estimate.</p></div>"), ti = /* @__PURE__ */ X("<div class=\"estimate-builder pb-16\"><div class=\"sticky top-0 z-10 bg-slate-900 text-white px-4 py-3 flex items-center justify-between border-b border-slate-700\"><div><h1 class=\"text-lg font-semibold\"> </h1> <span class=\"text-xs text-slate-400\">Estimate Builder</span></div> <div class=\"flex items-center gap-3\"><span class=\"text-xs px-2 py-1 rounded bg-slate-700 text-slate-300 uppercase tracking-wide\"> </span></div></div> <div class=\"bg-slate-50 border-b border-slate-200 px-4 py-2\"><div class=\"flex items-center gap-4 text-sm\"><span class=\"font-medium text-slate-600\">Global Markup:</span> <span class=\"px-2 py-0.5 rounded bg-blue-50 text-blue-700 font-mono text-xs\"> </span> <span class=\"px-2 py-0.5 rounded bg-amber-50 text-amber-700 font-mono text-xs\"> </span> <span class=\"px-2 py-0.5 rounded bg-purple-50 text-purple-700 font-mono text-xs\"> </span> <span class=\"px-2 py-0.5 rounded bg-green-50 text-green-700 font-mono text-xs\"> </span> <span class=\"px-2 py-0.5 rounded bg-slate-100 text-slate-600 font-mono text-xs\"> </span></div></div> <div class=\"p-4\"><!></div> <!></div>");
function ni(e, t) {
	Me(t, !0);
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
	tn(() => {
		t.projectId && a();
	});
	var o = $n(), s = Wt(o), c = (e) => {
		Z(e, Qr());
	}, l = (e) => {
		var t = $r(), n = L(t), i = L(n, !0);
		O(n), O(t), z(() => Q(i, Y(r))), Z(e, t);
	}, u = (e) => {
		var t = ti(), r = L(t), i = L(r), a = L(i), o = L(a, !0);
		O(a), Ce(2), O(i);
		var s = R(i, 2), c = L(s), l = L(c, !0);
		O(c), O(s), O(r);
		var u = R(r, 2), d = L(u), f = R(L(d), 2), p = L(f);
		O(f);
		var m = R(f, 2), h = L(m);
		O(m);
		var g = R(m, 2), _ = L(g);
		O(g);
		var v = R(g, 2), y = L(v);
		O(v);
		var b = R(v, 2), x = L(b);
		O(b), O(d), O(u);
		var S = R(u, 2), ee = L(S), C = (e) => {
			Z(e, ei());
		}, te = (e) => {
			var t = $n();
			lr(Wt(t), 17, () => Y(n).sections, (e) => e.id, (e, t) => {
				qr(e, {
					get section() {
						return Y(t);
					},
					get globals() {
						return Y(n).globals;
					}
				});
			}), Z(e, t);
		};
		$(ee, (e) => {
			Y(n).sections.length === 0 ? e(C) : e(te, -1);
		}), O(S), Zr(R(S, 2), { get estimate() {
			return Y(n);
		} }), O(t), z((e, t, r, i, a) => {
			Q(o, Y(n).project.name), Q(l, Y(n).project.status), Q(p, `Materials ${e ?? ""}`), Q(h, `Labor ${t ?? ""}`), Q(_, `Equipment ${r ?? ""}`), Q(y, `Subs ${i ?? ""}`), Q(x, `Other ${a ?? ""}`);
		}, [
			() => kr(Y(n).globals.materials_markup),
			() => kr(Y(n).globals.labor_markup),
			() => kr(Y(n).globals.equipment_markup),
			() => kr(Y(n).globals.subs_markup),
			() => kr(Y(n).globals.other_markup)
		]), Z(e, t);
	};
	$(s, (e) => {
		Y(i) ? e(c) : Y(r) ? e(l, 1) : Y(n) && e(u, 2);
	}), Z(e, o), Ne();
}
//#endregion
//#region src/main.js
var ri = document.getElementById("estimate-root");
if (ri) {
	let e = ri.dataset.projectId;
	er(ni, {
		target: ri,
		props: { projectId: e }
	});
}
//#endregion
