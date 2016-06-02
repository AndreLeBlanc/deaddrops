# Deadrop

**OSPP (1DT096) 2016 - Grupp 03**

Det här projektet ämnar att underlätta temporär filöverföring mellan två eller flera parter. En användare kan ladda upp en eller flera filer samt ange livslängd i form av dagar eller antal nedladdningar. Därefter så sparas de under en hashad URL och väntar på besökare. Servern som Deadrop körs på är skriven i Go och front-end i Javascript.

## Installera

### Back-end

##### Steg 1, installera Go
För att installera Go och se till så att ens GOROOT och GOPATH är rätt konfigurerade, rekommenderas det att följa
[Golangs officiella installationsguide](https://golang.org/doc/install).

##### Steg 2, klona projektet
Se till att stå i `$GOPATH/src` och skriv sen:

`git clone https://github.com/uu-it-teaching/ospp-2016-group-03/`

Därefter gå in i den nya katalogen och skriv:

`make dep` TODO: Fix this rule in the Makefile (make dependencies), change name of project folder to Deadrop and install database, etc.

### Front-end

##### Steg 1, installera npm
För att installera npm se till och föja deras [officiella installtionsguide](https://docs.npmjs.com/getting-started/installing-node)

##### Steg 2, installera gulp
Installera gulp som en global variabel. Om du kör linux använd 'sudo npm install gulp -g'. Man kan även följa denna [guide](https://github.com/gulpjs/gulp/blob/master/docs/getting-started.md)

##### Steg 3, installera webpack & webpack-dev-server
Nu kommer vi behöva installera webpack och webpack-dev-server globalt. Vi gör detta på en linux genom dessa kommandon:
'npm install webpack -g'
'npm install webpack-dev-server -g'

Alternativt följ [detta](https://webpack.github.io/docs/tutorials/getting-started/)

##### Steg 4, installera alla dependencies
Nu kör vi endast 'npm install' kommandot. Om allt går bra ska man nu vara redo och kompilera.

## Kompilera

### Back-end
Bygg hela projektet med:

`make build`

### Front-end

##### Alternativ 1, kompilera en gång
För att kompilera filerna endast en gång kan man köra dessa kommandon:
'npm run build_gulp'  - Detta bygger alla css filer
'npm run build_webpack' - Detta bygger alla javascript filer

Nu ska man kunna visa sidan, dock så är alternativ 2 både smidigare och snabbare.

##### Alternativ 2, aktiv kompilering + server
Vi börjar genom att se till att gulp lyssnar på ändringar i css filerna:
'gulp -watch &' - & lades till för att köra processen i bakgrunden
Sedan så kör vi igång vår dev-server som kommer att automatiskt bygga nya javascript filer vid behov:
'npm run dev' - nu körs servern på http:localhost:8080

Så lätt var det!

## Testa

### Back-end
Kör alla tester

`make test`

Kör alla tester med verbos flagga, om något test misslyckas skrivs det ut vart det blir fel:

`make testv`

### Front-end

## Starta systemet

### Back-end
För att starta servern kör:

`make run` TODO: Fix this rule

### Front-end

##### Precis som i "Alternativ 2, aktiv kompilering + server"
Vi börjar genom att se till att gulp lyssnar på ändringar i css filerna:
'gulp -watch &' - & lades till för att köra processen i bakgrunden
Sedan så kör vi igång vår dev-server som kommer att automatiskt bygga nya javascript filer vid behov:
'npm run dev' - nu körs servern på http:localhost:8080

## Struktur

Projektet består av följande kataloger.

### doc

projektrapporter och andra viktiga dokument.

### meta

- Presentation av gruppens medlemmar.
- Gruppkontrakt.
- Projektdagböcker.
- Reflektioner på gruppens arbete.

### api, database, server

All källkod.
