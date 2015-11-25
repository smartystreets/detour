// package detour offers an alternate, MVC-based, approach to HTTP applications.
// Rather than writing traditional http handlers you define input models that
// have optional Bind() and Validate() methods and which can be passed into
// methods on structs which return a Renderer. Each of these concepts is glued
// together by the ActionHandler struct via the New() function. See the example
// folder for a complete example.
package detour
