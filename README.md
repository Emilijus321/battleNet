BattleNet – Filmų Apžvalgų Platforma
Moderni, pilnai funkcionali internetinė aplikacija filmų atradimui, apžvalgoms ir asmeninės „watchlist“ valdymui. Sukurta naudojant Go, PostgreSQL, HTMX ir Templ sklandžiam vartotojo patyrimui.

✨ Funkcijos
👤 Vartotojo Funkcijos:
Saugus autentifikavimas: Bcrypt slaptažodžių hash’inimas su sesijų valdymu
Asmeninė „Watchlist“: Pridėkite ir valdykite mėgstamus filmus
Filmų apžvalgos: Įvertinkite (1–10) ir rašykite detalias apžvalgas su spoilerio įspėjimais
Profilio valdymas: Atnaujinkite asmeninę informaciją ir keiskite slaptažodį

👨‍💼 Administratoriaus Funkcijos:
Pilnas filmų valdymas: CRUD operacijos su filmais
TMDB integracija: Vienu paspaudimu importuokite filmus iš TMDB

👮 Moderatoriaus Funkcijos:
Vartotojų valdymas: Peržiūrėkite ir tvarkykite visus vartotojus
Rolės priskyrimas: Atnaujinkite vartotojų teises
Paskyrų moderavimas: Deaktyvuokite problematiškus vartotojus

🛠️ Technologijų Stack
Backend:
Programavimo kalba: Go 1.21+,
Web framework: Chi Router v5,
Duomenų bazė: PostgreSQL 14+ su pgx driver,
Migracijos: Goose,
Sesijų valdymas: SCS,
Slaptažodžių hash’inimas: bcrypt

Frontend
Templating: Templ – type-safe Go šablonai,
Interaktyvumas: HTMX 1.9.10,
Stilius: Custom CSS, responsive dizainas,
Ikonos: Unicode emoji

Išoriniai API:
TMDB API: Filmų duomenų bazė ir metaduomenys

1. Prieš pradedant
Įsitikinkite, kad įdiegta:
Go: 1.21 ar naujesnė,
PostgreSQL: 14 ar naujesnė,
Goose: Migrations įrankis

go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/a-h/templ/cmd/templ@latest

2. Klonuokite repozitoriją
git clone https://github.com/Emilijus321/battleNet.git
cd battleNet

3. Įdiekite priklausomybes
go mod download

4. Sukonfigūruokite aplinkos kintamuosius
Sukurkite .env failą:

# Duomenų bazės konfiguracija
DATABASE_URL=postgres://postgres:password@localhost:5432/movieapp?sslmode=disable

# Serverio konfiguracija
PORT=8080
ENVIRONMENT=development

# Saugumas
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# TMDB API
TMDB_API_KEY=your-tmdb-api-key-here
TMDB_BASE_URL=https://api.themoviedb.org/3
//TMDB API raktas: Sukurkite nemokamą TMDB paskyrą

5. Duomenų bazės nustatymas
createdb movieapp
goose -dir migrations postgres "postgres://postgres:password@localhost:5432/movieapp?sslmode=disable" up
goose -dir migrations postgres "postgres://postgres:password@localhost:5432/movieapp?sslmode=disable" status

6. Templ šablonų generavimas
templ generate

7. Paleiskite aplikaciją
go run main.go

Serveris bus pasiekiamas adresu: http://localhost:8080

