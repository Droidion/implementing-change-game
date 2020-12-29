drop table if exists roles cascade;

create table roles
(
    id   serial not null
        constraint roles_pk
            primary key,
    name text   not null
);

comment
on table roles is 'User roles';

create
unique index if not exists roles_id_uindex
    on roles (id);

INSERT INTO public.roles (id, name)
VALUES (1, 'Team');
INSERT INTO public.roles (id, name)
VALUES (2, 'Admin');

drop table if exists users;
create table users
(
    id       serial  not null
        constraint users_pk
            primary key,
    login    text    not null,
    password text    not null,
    role_id  integer not null
        constraint users_roles_id_fk
            references roles
            on update cascade on delete cascade
);

create
unique index if not exists users_id_uindex
    on users (id);

INSERT INTO public.users (id, login, password, role_id)
VALUES (1, 'admin', 'argon2id$19$65536$3$2$TONj4f02Nhy8JNM3tR1P+w$R55R6Z7J2xC6VZhxWcZcwizFrN+CK7BCPkiySo+s+8s', 2);