# Sistēmas koncepcijas prezentācija — satura plāns

Fails: `docs/presentation/sistemas_koncepcijas_prezentacija.pptx` (12 slaidi, 16:9, latviešu valoda).
Blakus ir PDF eksports `sistemas_koncepcijas_prezentacija.pdf` pārbaudei.
Prezentācijas mērķa laiks: ~4 min 50 s (limits 5 min). Katram slaidam ir runas
piezīmes PowerPoint piezīmju laukā — tās ir pavedieni, nevis lasāms teksts.

## Slaidi

| # | Slaids | Galvenā doma | Laiks |
| --- | --- | --- | --- |
| 1 | Titulslaids | Kas tiek izstrādāts, kam un kāpēc; autors, grupa DP3-1, 2026./2027. m.g. | 15 s |
| 2 | Saturs | Septiņas sadaļas | 10 s |
| 3 | Priekšmeta joma un problēma | Manuāls grafiks Excel tabulā; kļūdas atklājas tikai pēc sastādīšanas | 30 s |
| 4 | Tirgus izpēte | Excel vs. Planday / Deputy / Shiftbase; niša starp tiem | 35 s |
| 5 | Koncepcija: mērķis un mērķauditorija | Automatizēt grafiku; četras lietotāju lomas | 30 s |
| 6 | Galvenā funkcionalitāte | 6 funkciju bloki | 30 s |
| 7 | Kā notiek ģenerēšana | Iestatījumi → norma → sadale → pārbaudes → rezultāts | 30 s |
| 8 | Lomas un piekļuves tiesības | Hierarhija: reģistrēts lietotājs → darbinieks, vecākais darbinieks un sistēmas administrators; uzņēmuma administrators → vecākais darbinieks | 25 s |
| 9 | Lietošanas gadījumu diagramma | Attēls no `docs/media/usecase.png` + «Kā lasīt diagrammu» (ieskaitot aktoru vispārināšanu) | 15 s |
| 10 | SVID analīze | S/W pret funkcionalitāti, O/T pret tirgu; algoritma validācija kā pieņēmums | 30 s |
| 11 | Secinājumi | Kāpēc nepieciešama, ieguvumi, ierobežojumi un attīstība | 30 s |
| 12 | Informācijas avoti | 3 komerciālie risinājumi, Darba likums, saruna ar kasieri | 10 s |

## Pirms iesniegšanas jāpārbauda

2. **12. slaids** — precizēt URL un piekļuves datumus:
   - Planday: `https://www.planday.com` (piekļuves datums: ...)
   - Deputy: `https://www.deputy.com` (piekļuves datums: ...)
   - Shiftbase: `https://www.shiftbase.com` (piekļuves datums: ...)
   - Latvijas Republikas Darba likums — norādīt konkrētos pantus (darba laiks,
     atpūtas laiks), ja tie tiek minēti darbā.
   - Saruna ar mazumtirdzniecības kasieri vienā no lielākajiem mazumtirdzniecības
     veikaliem Latvijā — norādīt datumu. (Prezentācijā uzņēmuma nosaukums nav
     minēts, lai tas atbilstu darba dokumentācijai; ja drīkst, var norādīt IKEA.)

## Jautājumi, uz kuriem jābūt gatavam

- **Kāpēc neizmantot Planday/Deputy?** Maksa par lietotāju, pilnapjoma apjoms,
  automātiskā sastādīšana tikai augstākajā tarifā/agrīnā piekļuvē, Latvijas
  noteikumu atbalsts nav apstiprināts.
- **Kā nodrošina taisnīgumu?** Prioritāte darbiniekam ar mazāk ieplānotajām
  stundām; visiem redzams grafiks; apstiprināšanas un labojumu vēsture.
- **Kā ar Darba likumu?** Maksimālais maiņas garums un minimālais atpūtas laiks
  ir iestatāmi parametri, tāpēc tos var pielāgot vietējām prasībām.
- **Ko sistēma nedara?** Nav darba laika uzskaites; vienkāršots maiņu modelis
  (nakts/sadalītās maiņas, vecākie darbinieki netiek plānoti); nav mobilās
  lietotnes un algu integrācijas.
- **Kāpēc SVID vājās puses nav koncepcijas kļūdas?** Tās izriet no apzināta
  prototipa apjoma; attīstības virzieni ir prioritizēti 2.2. tabulā.
- **Kā pārliecināties, ka algoritms strādā reālos apstākļos?** Koncepcijā tā
  uzvedība sarežģītos grafikos vēl ir pieņēmums. Pirmais attīstības solis ir
  testēšana uz reāliem grafiku datiem ar pieņemšanas kritērijiem: segtā
  nodrošinājuma īpatsvars, stundu trūkums, normas un atpūtas pārkāpumu skaits.
  SVID iespējas (O) apzināti formulētas kā ārēji tirgus apstākļi, nevis kā
  iekšējie attīstības plāni.
- **Kā diagrammā lasāmas atkarības?** «include» – obligāta papilddarbība
  (piem., paroles ievade mainot paroli), «extend» – neobligāta (piem., filtru
  izvēle vēsturē); aktoru vispārināšana nozīmē tiesību pārmantošanu: darbinieks,
  vecākais darbinieks un sistēmas administrators manto reģistrētā lietotāja
  tiesības, bet uzņēmuma administrators – vecākā darbinieka tiesības.
- **Kāpēc sistēmas administrators nav uzņēmuma administratora apakšloma?**
  Diagrammā tas vispārina tieši reģistrēto lietotāju: tā tiesības ir uzņēmumu
  saraksts, uzņēmumu izveide un dzēšana, uzņēmuma administratora uzaicināšana un
  visu uzņēmumu darbību vēsture. Filiāļu, grafiku un darbinieku pārvaldību veic
  uzņēmuma administrators un vecākais darbinieks.
