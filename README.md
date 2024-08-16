# REST API using go
- Backend : Go
- Database : Postgres docker

<!-- #### Todo
- [x] Docker for Postgres
- [x] Initial db
- [ ] Router -->

### API
- <get>GET</get>`/api/v1/skills/:key` - get skill by key
- <get>GET</get>`/api/v1/skills` - get all skill
- <post>POST</post>`/api/v1/skills` - create skill
- <put>PUT</put>`/api/v1/skills/:key` - get skill by key
- <patch>PATCH</patch>`/api/v1/skills/name` - update name
- <patch>PATCH</patch>`/api/v1/skills/description` - update description
- <patch>PATCH</patch>`/api/v1/skills/logo` - update logo
- <patch>PATCH</patch>`/api/v1/skills/tags` - update tags
- <delete>DELETE</delete>`/api/v1/skills/:key` - delete skill by key

<style>
    get {
        background-color: #a4f084;
        color: #368228;
        font-weight: 500;
        padding: 0 0.5em;
        margin-right: 0.5em;
    }
    post {
        background-color: #fff878;
        color: #8c752a;
        font-weight: 500;
        padding: 0 0.5em;
        margin-right: 0.5em;
    }
    put {
        background-color: #fac57f;
        color: #8c461b;
        font-weight: 500;
        padding: 0 0.5em;
        margin-right: 0.5em;
    }
    patch {
        background-color: #fac57f;
        color: #8c461b;
        font-weight: 500;
        padding: 0 0.5em;
        margin-right: 0.5em;
    }
    delete {
        background-color: #ffa1a1;
        color: #6e1313;
        font-weight: 500;
        padding: 0 0.5em;
        margin-right: 0.5em;
    }
</style>