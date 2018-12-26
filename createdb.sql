--
-- PostgreSQL database dump
--

-- Dumped from database version 11.1
-- Dumped by pg_dump version 11.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: hashes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.hashes (
    hid integer NOT NULL,
    hash text,
    "time" integer,
    sid integer
);


--
-- Name: hashes_hid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.hashes_hid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: hashes_hid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.hashes_hid_seq OWNED BY public.hashes.hid;


--
-- Name: songs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.songs (
    sid integer NOT NULL,
    song text
);


--
-- Name: songs_sid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.songs_sid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: songs_sid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.songs_sid_seq OWNED BY public.songs.sid;


--
-- Name: hashes hid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hashes ALTER COLUMN hid SET DEFAULT nextval('public.hashes_hid_seq'::regclass);


--
-- Name: songs sid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.songs ALTER COLUMN sid SET DEFAULT nextval('public.songs_sid_seq'::regclass);


--
-- Name: hashes hashes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hashes
    ADD CONSTRAINT hashes_pkey PRIMARY KEY (hid);


--
-- Name: songs songs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.songs
    ADD CONSTRAINT songs_pkey PRIMARY KEY (sid);


--
-- Name: songs songs_song_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.songs
    ADD CONSTRAINT songs_song_key UNIQUE (song);


--
-- Name: hashes sid; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hashes
    ADD CONSTRAINT sid FOREIGN KEY (sid) REFERENCES public.songs(sid) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--
