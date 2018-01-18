--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.2
-- Dumped by pg_dump version 9.6.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: members; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE members (
    id uuid NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    voter_id character varying(255) NOT NULL,
    unit_number character varying(255) NOT NULL,
    street_number character varying(255) DEFAULT ''::character varying NOT NULL,
    street_name character varying(255) DEFAULT ''::character varying NOT NULL,
    city character varying(255) DEFAULT ''::character varying NOT NULL,
    state character varying(255) DEFAULT ''::character varying NOT NULL,
    postal_code character varying(255) DEFAULT ''::character varying NOT NULL,
    home_phone character varying(255) DEFAULT ''::character varying NOT NULL,
    cell_phone character varying(255) DEFAULT ''::character varying NOT NULL,
    recruiter character varying(255) DEFAULT ''::character varying NOT NULL,
    recruiter_phone character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    poll_id uuid NOT NULL,
    supporter boolean DEFAULT false NOT NULL,
    voted boolean DEFAULT false NOT NULL
);


ALTER TABLE members OWNER TO postgres;

--
-- Name: polls; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE polls (
    id uuid NOT NULL,
    name character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE polls OWNER TO postgres;

--
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE roles (
    id uuid NOT NULL,
    name character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE roles OWNER TO postgres;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE schema_migration (
    version character varying(255) NOT NULL
);


ALTER TABLE schema_migration OWNER TO postgres;

--
-- Name: user_roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE user_roles (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    role_id uuid NOT NULL
);


ALTER TABLE user_roles OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE users (
    id uuid NOT NULL,
    email character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    auth_id character varying(255) NOT NULL
);


ALTER TABLE users OWNER TO postgres;

--
-- Name: members members_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY members
    ADD CONSTRAINT members_pkey PRIMARY KEY (id);


--
-- Name: polls polls_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY polls
    ADD CONSTRAINT polls_pkey PRIMARY KEY (id);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: user_roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: polls_name_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX polls_name_idx ON polls USING btree (name);


--
-- Name: roles_name_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX roles_name_idx ON roles USING btree (name);


--
-- Name: user_roles_user_id_role_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_roles_user_id_role_id_idx ON user_roles USING btree (user_id, role_id);


--
-- Name: users_email_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX users_email_idx ON users USING btree (email);


--
-- Name: version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX version_idx ON schema_migration USING btree (version);


--
-- Name: members members_poll_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY members
    ADD CONSTRAINT members_poll_id_fkey FOREIGN KEY (poll_id) REFERENCES polls(id);


--
-- Name: user_roles user_roles_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY user_roles
    ADD CONSTRAINT user_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES roles(id);


--
-- Name: user_roles user_roles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY user_roles
    ADD CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);


--
-- PostgreSQL database dump complete
--

