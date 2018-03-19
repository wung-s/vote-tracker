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
-- Name: tiger; Type: SCHEMA; Schema: -; Owner: Wung
--

CREATE SCHEMA tiger;


ALTER SCHEMA tiger OWNER TO "Wung";

--
-- Name: tiger_data; Type: SCHEMA; Schema: -; Owner: Wung
--

CREATE SCHEMA tiger_data;


ALTER SCHEMA tiger_data OWNER TO "Wung";

--
-- Name: topology; Type: SCHEMA; Schema: -; Owner: Wung
--

CREATE SCHEMA topology;


ALTER SCHEMA topology OWNER TO "Wung";

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: address_standardizer; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS address_standardizer WITH SCHEMA public;


--
-- Name: EXTENSION address_standardizer; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION address_standardizer IS 'Used to parse an address into constituent elements. Generally used to support geocoding address normalization step.';


--
-- Name: address_standardizer_data_us; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS address_standardizer_data_us WITH SCHEMA public;


--
-- Name: EXTENSION address_standardizer_data_us; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION address_standardizer_data_us IS 'Address Standardizer US dataset example';


--
-- Name: fuzzystrmatch; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS fuzzystrmatch WITH SCHEMA public;


--
-- Name: EXTENSION fuzzystrmatch; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION fuzzystrmatch IS 'determine similarities and distance between strings';


--
-- Name: postgis; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA public;


--
-- Name: EXTENSION postgis; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION postgis IS 'PostGIS geometry, geography, and raster spatial types and functions';


--
-- Name: postgis_sfcgal; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS postgis_sfcgal WITH SCHEMA public;


--
-- Name: EXTENSION postgis_sfcgal; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION postgis_sfcgal IS 'PostGIS SFCGAL functions';


--
-- Name: postgis_tiger_geocoder; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS postgis_tiger_geocoder WITH SCHEMA tiger;


--
-- Name: EXTENSION postgis_tiger_geocoder; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION postgis_tiger_geocoder IS 'PostGIS tiger geocoder and reverse geocoder';


--
-- Name: postgis_topology; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS postgis_topology WITH SCHEMA topology;


--
-- Name: EXTENSION postgis_topology; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION postgis_topology IS 'PostGIS topology spatial types and functions';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: dispositions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE dispositions (
    id uuid NOT NULL,
    intention character varying(255) DEFAULT ''::character varying NOT NULL,
    contact_type character varying(255) DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    member_id uuid NOT NULL
);


ALTER TABLE dispositions OWNER TO postgres;

--
-- Name: electoral_districts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE electoral_districts (
    id uuid NOT NULL,
    name character varying(255) DEFAULT ''::character varying NOT NULL,
    edid integer,
    shape_area numeric,
    shape_length numeric,
    geom geometry(Polygon,4326) NOT NULL
);


ALTER TABLE electoral_districts OWNER TO postgres;

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
    recruiter_phone character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    poll_id uuid NOT NULL,
    supporter boolean DEFAULT false NOT NULL,
    voted boolean DEFAULT false NOT NULL,
    recruiter character varying(255) DEFAULT ''::character varying NOT NULL,
    recruiter_id uuid NOT NULL,
    latlng geometry(Point,4326)
);


ALTER TABLE members OWNER TO postgres;

--
-- Name: members_view; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW members_view AS
 SELECT members.id,
    members.first_name,
    members.last_name,
    members.voter_id,
    members.unit_number,
    members.street_number,
    members.street_name,
    members.city,
    members.state,
    members.postal_code,
    members.home_phone,
    members.cell_phone,
    members.recruiter_phone,
    members.created_at,
    members.updated_at,
    members.poll_id,
    members.supporter,
    members.voted,
    members.recruiter,
    members.recruiter_id,
    members.latlng,
    concat_ws(' '::text, (members.unit_number)::text, (members.street_number)::text, (members.street_name)::text) AS address
   FROM members;


ALTER TABLE members_view OWNER TO postgres;

--
-- Name: polling_divisions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE polling_divisions (
    id uuid NOT NULL,
    edid integer,
    no integer,
    shape_area numeric,
    shape_length numeric
);


ALTER TABLE polling_divisions OWNER TO postgres;

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
-- Name: recruiters; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE recruiters (
    id uuid NOT NULL,
    name character varying(255) DEFAULT ''::character varying NOT NULL,
    phone_no character varying(255) DEFAULT ''::character varying NOT NULL,
    invited boolean DEFAULT false NOT NULL,
    notification_enabled boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE recruiters OWNER TO postgres;

--
-- Name: ride_requests; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE ride_requests (
    id uuid NOT NULL,
    address text DEFAULT ''::text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    member_id uuid NOT NULL
);


ALTER TABLE ride_requests OWNER TO postgres;

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
    auth_id character varying(255) NOT NULL,
    phone_no character varying(255) DEFAULT ''::character varying NOT NULL,
    poll_id uuid,
    invited boolean
);


ALTER TABLE users OWNER TO postgres;

--
-- Name: dispositions dispositions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY dispositions
    ADD CONSTRAINT dispositions_pkey PRIMARY KEY (id);


--
-- Name: electoral_districts electoral_districts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY electoral_districts
    ADD CONSTRAINT electoral_districts_pkey PRIMARY KEY (id);


--
-- Name: members members_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY members
    ADD CONSTRAINT members_pkey PRIMARY KEY (id);


--
-- Name: polling_divisions polling_divisions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY polling_divisions
    ADD CONSTRAINT polling_divisions_pkey PRIMARY KEY (id);


--
-- Name: polls polls_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY polls
    ADD CONSTRAINT polls_pkey PRIMARY KEY (id);


--
-- Name: recruiters recruiters_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY recruiters
    ADD CONSTRAINT recruiters_pkey PRIMARY KEY (id);


--
-- Name: ride_requests ride_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY ride_requests
    ADD CONSTRAINT ride_requests_pkey PRIMARY KEY (id);


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
-- Name: electoral_districts_edid_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX electoral_districts_edid_idx ON electoral_districts USING btree (edid);


--
-- Name: electoral_districts_geom_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX electoral_districts_geom_idx ON electoral_districts USING gist (geom);


--
-- Name: polls_name_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX polls_name_idx ON polls USING btree (name);


--
-- Name: recruiters_phone_no_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX recruiters_phone_no_idx ON recruiters USING btree (phone_no);


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
-- Name: dispositions dispositions_member_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY dispositions
    ADD CONSTRAINT dispositions_member_id_fkey FOREIGN KEY (member_id) REFERENCES members(id);


--
-- Name: members members_poll_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY members
    ADD CONSTRAINT members_poll_id_fkey FOREIGN KEY (poll_id) REFERENCES polls(id);


--
-- Name: members members_recruiter_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY members
    ADD CONSTRAINT members_recruiter_id_fkey FOREIGN KEY (recruiter_id) REFERENCES recruiters(id);


--
-- Name: ride_requests ride_requests_member_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY ride_requests
    ADD CONSTRAINT ride_requests_member_id_fkey FOREIGN KEY (member_id) REFERENCES members(id);


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
-- Name: users users_poll_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_poll_id_fkey FOREIGN KEY (poll_id) REFERENCES polls(id);


--
-- PostgreSQL database dump complete
--

