// package detour offers an alternate, MVC-based, approach to HTTP applications.
// Rather than writing traditional http.Handlers you define input models that
// have optional Bind(), Sanitize(), and Validate() methods and which can be passed into
// methods on structs which return a Renderer. Each of these concepts is glued
// together by the library's ActionHandler struct via the New() function. See the example
// folder for a complete example.
// Requires Go 1.7+
package detour
