---
description: >-
  On this page a brief explanation of the need for essentia and the ideas behind
  it are explained.
---

# Philosophy

The simplicity of Go is both a gift and a "curse". It eliminates a lot of overhead while also providing quite high level standard libraries (eg. net/http) to developers. The "curse" lies in the fact that Go then also becomes quite repetitive and you basically have to create your own API framework/library.

This is also what _essentia_ tries to resolve while keeping performance in your hands. The core idea of _essentia_ is to have a framework in which only higher level API framework functionalities (ie. handlers) and middleware (ie. router) are implemented. These are the layers of an API that often remain the same across projects and where performance can't be improved.

Lower level API framework functionalities are thus not implemented but are called in these implemented higher level functionalities. You can see _essentia_ as your standard web/API framework without ORM. This because ORMs are quite useful for simpler use cases but are harder to overcome than actually implement for more complicated use cases.

This philosophy is summarized in the name, _essentia_. It allows developers to focus on the **essence** of their API: achieving optimal performance without compromises.



