--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4
-- Dumped by pg_dump version 17.4

-- Started on 2025-04-26 17:22:48

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 7 (class 2615 OID 16386)
-- Name: pgagent; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA pgagent;


ALTER SCHEMA pgagent OWNER TO postgres;

--
-- TOC entry 5056 (class 0 OID 0)
-- Dependencies: 7
-- Name: SCHEMA pgagent; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA pgagent IS 'pgAgent system tables';


--
-- TOC entry 2 (class 3079 OID 16387)
-- Name: pgagent; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgagent WITH SCHEMA pgagent;


--
-- TOC entry 5057 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION pgagent; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgagent IS 'A PostgreSQL job scheduler';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 239 (class 1259 OID 16562)
-- Name: clothing_advice; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.clothing_advice (
    id integer NOT NULL,
    comments text NOT NULL,
    style text NOT NULL,
    temp_max smallint,
    temp_min smallint,
    conditions text[],
    max_speed real,
    season text,
    accessories text
);


ALTER TABLE public.clothing_advice OWNER TO postgres;

--
-- TOC entry 238 (class 1259 OID 16561)
-- Name: clothing_advice_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.clothing_advice_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.clothing_advice_id_seq OWNER TO postgres;

--
-- TOC entry 5058 (class 0 OID 0)
-- Dependencies: 238
-- Name: clothing_advice_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.clothing_advice_id_seq OWNED BY public.clothing_advice.id;


--
-- TOC entry 242 (class 1259 OID 16634)
-- Name: conditions_advice; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.conditions_advice (
    conditions text,
    conditions_comments text,
    temp_min smallint,
    temp_max smallint,
    id integer NOT NULL
);


ALTER TABLE public.conditions_advice OWNER TO postgres;

--
-- TOC entry 245 (class 1259 OID 16689)
-- Name: conditions_advice_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.conditions_advice_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.conditions_advice_id_seq OWNER TO postgres;

--
-- TOC entry 5059 (class 0 OID 0)
-- Dependencies: 245
-- Name: conditions_advice_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.conditions_advice_id_seq OWNED BY public.conditions_advice.id;


--
-- TOC entry 241 (class 1259 OID 16571)
-- Name: mems; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mems (
    id integer NOT NULL,
    temp smallint NOT NULL,
    mems text NOT NULL
);


ALTER TABLE public.mems OWNER TO postgres;

--
-- TOC entry 240 (class 1259 OID 16570)
-- Name: mems_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.mems_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.mems_id_seq OWNER TO postgres;

--
-- TOC entry 5060 (class 0 OID 0)
-- Dependencies: 240
-- Name: mems_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.mems_id_seq OWNED BY public.mems.id;


--
-- TOC entry 243 (class 1259 OID 16639)
-- Name: pressure_advice; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pressure_advice (
    min_pressure smallint,
    max_pressure smallint,
    pressure_comments text
);


ALTER TABLE public.pressure_advice OWNER TO postgres;

--
-- TOC entry 237 (class 1259 OID 16551)
-- Name: weather_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.weather_history (
    id integer NOT NULL,
    city text NOT NULL,
    temp smallint NOT NULL,
    conditions text,
    pressure smallint,
    wind_speed real,
    date timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.weather_history OWNER TO postgres;

--
-- TOC entry 236 (class 1259 OID 16550)
-- Name: weather_history_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.weather_history_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.weather_history_id_seq OWNER TO postgres;

--
-- TOC entry 5061 (class 0 OID 0)
-- Dependencies: 236
-- Name: weather_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.weather_history_id_seq OWNED BY public.weather_history.id;


--
-- TOC entry 244 (class 1259 OID 16644)
-- Name: wind_advice; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.wind_advice (
    min_speed real,
    max_speed real,
    wind_comments text,
    id integer NOT NULL
);


ALTER TABLE public.wind_advice OWNER TO postgres;

--
-- TOC entry 246 (class 1259 OID 16699)
-- Name: wind_advice_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.wind_advice_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.wind_advice_id_seq OWNER TO postgres;

--
-- TOC entry 5062 (class 0 OID 0)
-- Dependencies: 246
-- Name: wind_advice_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.wind_advice_id_seq OWNED BY public.wind_advice.id;


--
-- TOC entry 4847 (class 2604 OID 16565)
-- Name: clothing_advice id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clothing_advice ALTER COLUMN id SET DEFAULT nextval('public.clothing_advice_id_seq'::regclass);


--
-- TOC entry 4849 (class 2604 OID 16690)
-- Name: conditions_advice id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.conditions_advice ALTER COLUMN id SET DEFAULT nextval('public.conditions_advice_id_seq'::regclass);


--
-- TOC entry 4848 (class 2604 OID 16574)
-- Name: mems id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mems ALTER COLUMN id SET DEFAULT nextval('public.mems_id_seq'::regclass);


--
-- TOC entry 4845 (class 2604 OID 16554)
-- Name: weather_history id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.weather_history ALTER COLUMN id SET DEFAULT nextval('public.weather_history_id_seq'::regclass);


--
-- TOC entry 4850 (class 2604 OID 16700)
-- Name: wind_advice id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wind_advice ALTER COLUMN id SET DEFAULT nextval('public.wind_advice_id_seq'::regclass);


--
-- TOC entry 4807 (class 0 OID 16388)
-- Dependencies: 221
-- Data for Name: pga_jobagent; Type: TABLE DATA; Schema: pgagent; Owner: postgres
--

COPY pgagent.pga_jobagent (jagpid, jaglogintime, jagstation) FROM stdin;
19172	2025-04-17 13:17:00.774451+03	DESKTOP-S2R1O15
\.


--
-- TOC entry 4808 (class 0 OID 16397)
-- Dependencies: 223
-- Data for Name: pga_jobclass; Type: TABLE DATA; Schema: pgagent; Owner: postgres
--

COPY pgagent.pga_jobclass (jclid, jclname) FROM stdin;
\.


--
-- TOC entry 4809 (class 0 OID 16407)
-- Dependencies: 225
-- Data for Name: pga_job; Type: TABLE DATA; Schema: pgagent; Owner: postgres
--

COPY pgagent.pga_job (jobid, jobjclid, jobname, jobdesc, jobhostagent, jobenabled, jobcreated, jobchanged, jobagentid, jobnextrun, joblastrun) FROM stdin;
\.


--
-- TOC entry 4811 (class 0 OID 16455)
-- Dependencies: 229
-- Data for Name: pga_schedule; Type: TABLE DATA; Schema: pgagent; Owner: postgres
--

COPY pgagent.pga_schedule (jscid, jscjobid, jscname, jscdesc, jscenabled, jscstart, jscend, jscminutes, jschours, jscweekdays, jscmonthdays, jscmonths) FROM stdin;
\.


--
-- TOC entry 4812 (class 0 OID 16483)
-- Dependencies: 231
-- Data for Name: pga_exception; Type: TABLE DATA; Schema: pgagent; Owner: postgres
--

COPY pgagent.pga_exception (jexid, jexscid, jexdate, jextime) FROM stdin;
\.


--
-- TOC entry 4813 (class 0 OID 16497)
-- Dependencies: 233
-- Data for Name: pga_joblog; Type: TABLE DATA; Schema: pgagent; Owner: postgres
--

COPY pgagent.pga_joblog (jlgid, jlgjobid, jlgstatus, jlgstart, jlgduration) FROM stdin;
\.


--
-- TOC entry 4810 (class 0 OID 16431)
-- Dependencies: 227
-- Data for Name: pga_jobstep; Type: TABLE DATA; Schema: pgagent; Owner: postgres
--

COPY pgagent.pga_jobstep (jstid, jstjobid, jstname, jstdesc, jstenabled, jstkind, jstcode, jstconnstr, jstdbname, jstonerror, jscnextrun) FROM stdin;
\.


--
-- TOC entry 4814 (class 0 OID 16513)
-- Dependencies: 235
-- Data for Name: pga_jobsteplog; Type: TABLE DATA; Schema: pgagent; Owner: postgres
--

COPY pgagent.pga_jobsteplog (jslid, jsljlgid, jsljstid, jslstatus, jslresult, jslstart, jslduration, jsloutput) FROM stdin;
\.


--
-- TOC entry 5043 (class 0 OID 16562)
-- Dependencies: 239
-- Data for Name: clothing_advice; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.clothing_advice (id, comments, style, temp_max, temp_min, conditions, max_speed, season, accessories) FROM stdin;
1	Верх: Термобельё + толстый свитер + пуховик/дутая куртка.\nНиз: Утеплённые брюки/джинсы + термоштаны.\nАксессуары: Шапка-ушанка, шарф, варежки, термоноски.\nОбувь: Зимние ботинки с мембраной (типа Sorel или Timberland).	Арктический исследователь	0	-15	{Clouds,Clear,Show}	15	Зима	Перчатки и шарф
2	Верх: Вязаный свитер оверсайз + длинное пальто (шерсть/кашемир).\nНиз: Узкие джинсы + термоколготки (если юбка).\nОбувь: Челси на толстой подошве или угги.	Скандинавский минимализм	0	-15	{Clouds,Clear,Show}	15	Зима	Перчатки и шарф
3	Только лучшее из "Охотник и рыбалов"	Зимняя рыбалка	-15	-30	{Clouds,Clear,Show,Thunderstorm}	100	Зима	Перчатки, шарф и бафф
8	Верх: Льняная рубашка + майка.\nНиз: Широкие льняные брюки или сарафан.\nОбувь: Сандалии Birkenstock или эспадрильи.	Богемный стиль	27	16	{Clouds,Clear}	12	Лето	Очки и кепку
4	Верх: Тренчкот + тонкий свитер + рубашка.\nНиз: Классические брюки или юбка-карандаш + колготки.\nОбувь: Чёрные кожаные ботинки Dr. Martens.	Британский стиль	10	0	{Clouds,Clear}	12	Осень	Лёгкий шарф
5	Верх: Фланелевая рубашка + толстовка + джинсовая куртка.\nНиз: Чиносы или джинсы с подвёрнутыми штанинами.\nОбувь: Кроссовки Nike Air Force или Vans.	Уютный хипстер	10	0	{Clouds,Clear}	12	Весна	Лёгкий шарф
6	Верх: Приталенный пиджак + водолазка.\nНиз: Юбка-миди или зауженные брюки.\nОбувь: Лоферы или балетки.	Французская классика	16	10	{Clouds,Clear}	12	Весна	Очки от ветра/пыли
7	Верх: Белая футболка + кожаная косуха.\nНиз: Джинсы mom-fit.\nОбувь: Белые кроссовки (Adidas Stan Smith).	Кэжуал-шик	16	10	{Clouds,Clear}	12	Осень	Очки от ветра/пыли
9	Верх: Кроп-топ + ветровка.\nНиз: Велосипедки или шорты-бермуды.\nОбувь: Кроссовки New Balance.	Спортивный гламур	27	16	{Clouds,Clear}	15	Лето	Очки и кепку
10	Купальник, очки и накаченный пресс	Пора на пляж	50	27	{Clear}	12	Лето	Очки, кепку и хорошее настроение
14	Лёгкий свитер + тренч + зонт. Мокасины или ботинки с высоким берцем. Непромокаемые брюки или джинсы.	Париж в ноябре	10	0	{Rain,Fog}	10	Осень	Зонт
16	Белая футболка + джинсовые шорты. Панама или бейсболка. Кеды или шлёпанцы.	Летние деньки	30	24	{Clear}	4	Лето	Кепку и очки
18	Под такую погоду стиля ещё не добавлено:( Жди обновления, совсем скоро мы это исправим.	Pop	100	-100	{Thunderstorm,Fog,Show,Rain,Drizzle,Clouds,Clear}	100	Любое	Ничего
17	На улице шторм, лучше только под крышей. Надень что-то лёгкое	Громовой ниндзя	50	-30	{Thunderstorm}	10	Любое	Громоотвод
13	Термобельё, флисовая кофта, горнолыжная куртка. Штаны на утеплителе + термогольфы.	Полярная экспедиция	-15	-30	{Snow}	10	Зима	Перчатки, шарф и бафф
15	Лонгслив или худи, укороченный бушлат. Джинсы или чиносы + кроссовки. 	Минимализм в движении	16	10	{Clear,Clouds}	6	Весна	Ничего
\.


--
-- TOC entry 5046 (class 0 OID 16634)
-- Dependencies: 242
-- Data for Name: conditions_advice; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.conditions_advice (conditions, conditions_comments, temp_min, temp_max, id) FROM stdin;
Clear	На улице солнце — не забудь очки и кепку! 😎	18	35	1
Clear	Идеальная погода для прогулки, наслаждайся! 🥳	18	35	2
Clear	Ясно, но не жарко — можно надеть лёгкую кофту, очки по желанию 🌤	0	17	3
Clear	Солнечно, но холодно — утеплись! Очки не обязательны 🧥	0	17	4
Clear	Солнечно, но холодно — утеплись! Очки не обязательны 🧥	-50	-1	23
Clear	На солнце может светить ярко — но холодно ❄️	-50	-1	24
Clouds	Пасмурно, но пока без дождя. Возьми кофту ☁️	18	35	5
Clouds	Небо затянуто — не лишним будет зонт на всякий 🌥️	18	35	6
Clouds	Хмуро снаружи, зато уютно внутри. Надень что-то мягкое, но тёплое ☁️	0	17	7
Clouds	Облака вокруг — на случай дождя возьми что-то с капюшоном 🌥️	0	17	8
Clouds	Морозно и тускло, при сильном ветре дополнительно защитить лицо баффом ☃️	-50	-1	25
Clouds	Холодно и мрачно, одевайся по погоде 🥼	-50	-1	26
Fog	На улице туман — надень что-то заметное, будь осторожен на дороге 🌫️	-50	50	9
Fog	Туман сгущается — захвати фонарик, если идёшь вечером 🌁	-50	50	10
Fog	Видимость плохая — не забудь светоотражающие элементы 🌫️	-50	50	11
Snow	Снег валит — пора утеплиться: шапка, шарф, варежки ☃️	-50	50	12
Snow	Скользко и снежно — подойдут ботинки с хорошей подошвой ⛄	-50	50	13
Snow	Снег кружится — отличное время для пуховика 🌨️	-50	50	14
Rain	На улице дождь — без зонта лучше не выходить 🌦️	-50	50	15
Rain	Дождливо — накинь дождевик и обувь по погоде 🌧️	-50	50	16
Rain	Мокро и сыро — не забудь зонт и непромокаемые ботинки ☔	-50	50	17
Drizzle	Моросящий дождик — зонт не помешает, но можно и с капюшоном 💦	-50	50	18
Drizzle	Слегка капает — надень что-то с водоотталкивающей тканью 💦	-50	50	19
Thunderstorm	Гремит гром — лучше остаться дома, если можно 🌩️	-50	50	20
Thunderstorm	Грозовая активность — будь осторожен с техникой и деревьями ⛈️	-50	50	21
Thunderstorm	Молнии сверкают — держись подальше от открытых пространств 🌩️	-50	50	22
\.


--
-- TOC entry 5045 (class 0 OID 16571)
-- Dependencies: 241
-- Data for Name: mems; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.mems (id, temp, mems) FROM stdin;
\.


--
-- TOC entry 5047 (class 0 OID 16639)
-- Dependencies: 243
-- Data for Name: pressure_advice; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.pressure_advice (min_pressure, max_pressure, pressure_comments) FROM stdin;
900	999	Низкое давление — возможна слабость и головные боли. Будь внимателен
900	999	Атмосферное давление понижено — сократи физнагрузки
1000	1025	Давление в пределах нормы — хорошее самочувствие ожидается
1000	1025	Давление в норме, за этот аспект можешь не переживать
1026	1100	Высокое давление — может чувствоваться дискомфорт в висках
1026	1100	Давление выше нормы — если метеозависим, будь осторожен
\.


--
-- TOC entry 5041 (class 0 OID 16551)
-- Dependencies: 237
-- Data for Name: weather_history; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.weather_history (id, city, temp, conditions, pressure, wind_speed, date) FROM stdin;
231	Воронеж	9	Clear	1008	8	2025-04-26 10:52:34.913185
232	Курск	6	Clouds	1012	5.54	2025-04-26 10:55:30.127603
233	Воронеж	9	Clear	1008	8	2025-04-26 10:57:39.326957
234	Курск	6	Clouds	1012	5.54	2025-04-26 11:00:30.291464
235	Воронеж	9	Clear	1008	8	2025-04-26 11:00:57.250032
236	Курск	6	Clouds	1012	5.54	2025-04-26 11:01:54.259811
237	Воронеж	9	Clear	1008	8	2025-04-26 11:02:07.117234
238	Воронеж	7	Clear	1009	8	2025-04-26 11:36:06.036285
239	Вава	18	Clouds	1011	2.32	2025-04-26 11:45:08.505846
240	Курск	4	Clouds	1012	5.87	2025-04-26 11:45:25.065321
241	Воронеж	2	Snow	1012	4	2025-04-26 15:11:59.033492
242	Воронеж	2	Snow	1012	4	2025-04-26 15:13:58.384223
243	Воронеж	2	Snow	1012	4	2025-04-26 15:14:40.677607
244	Воронеж	2	Snow	1012	4	2025-04-26 15:16:00.253152
245	Воронеж	2	Snow	1012	4	2025-04-26 15:16:43.908817
\.


--
-- TOC entry 5048 (class 0 OID 16644)
-- Dependencies: 244
-- Data for Name: wind_advice; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.wind_advice (min_speed, max_speed, wind_comments, id) FROM stdin;
0	5	Почти штиль, наслаждайся!	2
5	10	Если волосы важны — возьми головной убор 😉	4
20	100	Очень сильный ветер — избегай открытых пространств и зонтов!	7
20	100	Будет дуть во все стороны. Одежду лучше посерьёзнее, без "парусов"	8
0	5	Ветер почти отсутствует — температура ощущаемая равна фактической	1
5	10	Лёгкий ветерок — не обращай внимание	3
10	20	Сильный ветер — будь осторожен	5
10	20	Такой ветер — как бесплатный фен, учитывай, может быть холоднее	6
\.


--
-- TOC entry 5063 (class 0 OID 0)
-- Dependencies: 238
-- Name: clothing_advice_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.clothing_advice_id_seq', 19, true);


--
-- TOC entry 5064 (class 0 OID 0)
-- Dependencies: 245
-- Name: conditions_advice_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.conditions_advice_id_seq', 26, true);


--
-- TOC entry 5065 (class 0 OID 0)
-- Dependencies: 240
-- Name: mems_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.mems_id_seq', 1, false);


--
-- TOC entry 5066 (class 0 OID 0)
-- Dependencies: 236
-- Name: weather_history_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.weather_history_id_seq', 245, true);


--
-- TOC entry 5067 (class 0 OID 0)
-- Dependencies: 246
-- Name: wind_advice_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.wind_advice_id_seq', 8, true);


--
-- TOC entry 4888 (class 2606 OID 16569)
-- Name: clothing_advice clothing_advice_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clothing_advice
    ADD CONSTRAINT clothing_advice_pkey PRIMARY KEY (id);


--
-- TOC entry 4892 (class 2606 OID 16692)
-- Name: conditions_advice conditions_advice_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.conditions_advice
    ADD CONSTRAINT conditions_advice_pkey PRIMARY KEY (id);


--
-- TOC entry 4890 (class 2606 OID 16578)
-- Name: mems mems_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mems
    ADD CONSTRAINT mems_pkey PRIMARY KEY (id);


--
-- TOC entry 4886 (class 2606 OID 16559)
-- Name: weather_history weather_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.weather_history
    ADD CONSTRAINT weather_history_pkey PRIMARY KEY (id);


--
-- TOC entry 4894 (class 2606 OID 16702)
-- Name: wind_advice wind_advice_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wind_advice
    ADD CONSTRAINT wind_advice_pkey PRIMARY KEY (id);


-- Completed on 2025-04-26 17:22:48

--
-- PostgreSQL database dump complete
--

