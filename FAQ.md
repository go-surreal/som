# FAQ

*Disclaimer: Currently those are just questions I asked myself and wanted an answer before the initial public release.
In the future this section will be advanced by topics raised in issues or discussions.*

## Is this library performant?

The main goal of an ORM/OGM is not performance itself, but to make the life of a programmer 
using the ORM/OGM easier by abstracting the database access into a cleaner API.
Due to that abstraction, it cannot be faster than directly accessing the database layer with native queries.
So while performance is (and always will be) a concern of this library and will try to stay as fast as possible,
it is not the "top priority" here.

Still, if there are obvious performance bugs, please open an issue! 🙏

## Should I wrap th usage of this library in model repositories?

- https://cockneycoder.wordpress.com/2013/04/07/why-entity-framework-renders-the-repository-pattern-obsolete/
- https://lostechies.com/jimmybogard/2012/10/08/favor-query-objects-over-repositories/

### Why are maps not supported?

- With the schemaless database this would be possible.
- Currently, the focus is on structured and deterministic data.
- Might be added in the future though.

### Why does a filter like `where.User.Equal(userModel)` not exist?

- This would be an ambiguous case. Should it compare the whole object with all properties or only the ID?
- For this case it is better and more deterministic to just compare the ID explicitly.
- If - for whatever reason - it is required to check the fields, adding those filters one by one makes the purpose of the query clearer.
- Furthermore, we would need to find a way to circumvent a naming clash when the field of a model is named `Equal` (or other keywords).
- On the other hand, this feature is still open for debate. So if anyone can clarify the need for it, we might as well implement it at some point.
