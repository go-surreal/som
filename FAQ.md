# FAQ

## Is this library performant?

The main goal of an ORM/OGM is not performance itself, but to make the life of a programmer 
using the ORM/OGM easier by abstracting the database access into a cleaner API.
Due to that abstraction, it cannot be faster than directly accessing the database layer with native queries.
So while performance is (and always will be) a concern of this library and will try to stay as fast as possible,
it is not the "top priority" here.

Still, if there are obvious performance bugs, please open an issue! 🙏
