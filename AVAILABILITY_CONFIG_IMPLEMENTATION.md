# System Konfiguracji Dostępności Użytkownika - Dokumentacja Zmian

## Przegląd

Zaimplementowano nowy system matchingu oparty na **stałej konfiguracji dostępności użytkownika**.

### Poprzedni system:

- Użytkownicy określali dostępność **podczas/po** utworzeniu matcha (tabela `match_availabilities`)
- Match mógł być utworzony bez sprawdzania wspólnych dostępności

### Nowy system:

- Każdy użytkownik ma **własną konfigurację dostępności** (profil systemowy)
- Match powstaje **tylko wtedy**, gdy oboje użytkownicy:
  1. Mają utworzoną konfigurację dostępności
  2. Mają **co najmniej jedną wspólną dostępność** (dzień + pora)

## Zmiany w bazie danych

### Nowa tabela: `user_availability_configs`

```sql
CREATE TABLE IF NOT EXISTS user_availability_configs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    monday_morning BOOLEAN DEFAULT false,
    monday_afternoon BOOLEAN DEFAULT false,
    tuesday_morning BOOLEAN DEFAULT false,
    tuesday_afternoon BOOLEAN DEFAULT false,
    wednesday_morning BOOLEAN DEFAULT false,
    wednesday_afternoon BOOLEAN DEFAULT false,
    thursday_morning BOOLEAN DEFAULT false,
    thursday_afternoon BOOLEAN DEFAULT false,
    friday_morning BOOLEAN DEFAULT false,
    friday_afternoon BOOLEAN DEFAULT false,
    saturday_morning BOOLEAN DEFAULT false,
    saturday_afternoon BOOLEAN DEFAULT false,
    sunday_morning BOOLEAN DEFAULT false,
    sunday_afternoon BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

**Migracja:** `migrations/000015_create_user_availability_configs_table.up.sql`

## Nowe pliki

### 1. Model: `models/user_availability_config.go`

- `UserAvailabilityConfig` - główna struktura konfiguracji
- `CreateAvailabilityConfigInput` - input do tworzenia
- `UpdateAvailabilityConfigInput` - input do aktualizacji
- Funkcje pomocnicze:
  - `GetAvailableSlots()` - zwraca listę dostępnych slotów
  - `HasCommonAvailability(config1, config2)` - sprawdza wspólne dostępności
  - `GetCommonSlots(config1, config2)` - zwraca wspólne sloty

### 2. Repository: `repositories/user_availability_config_repository.go`

- `Create(config)` - tworzy nową konfigurację
- `FindByUserID(userID)` - pobiera konfigurację użytkownika
- `FindByUserIDs(userIDs)` - pobiera konfiguracje wielu użytkowników
- `Update(config)` - aktualizuje konfigurację
- `Delete(userID)` - usuwa konfigurację
- `Exists(userID)` - sprawdza istnienie konfiguracji

### 3. Service: `services/user_availability_config_service.go`

- `CreateConfig()` - tworzy konfigurację z walidacją
- `GetConfig()` - pobiera konfigurację
- `UpdateConfig()` - aktualizuje konfigurację z walidacją
- `DeleteConfig()` - usuwa konfigurację
- `HasConfig()` - sprawdza istnienie konfiguracji

**Walidacja:** Co najmniej jeden slot musi być wybrany

### 4. Handler: `handlers/user_availability_config_handler.go`

- `POST /api/availability-config` - tworzy konfigurację
- `GET /api/availability-config` - pobiera konfigurację
- `PUT /api/availability-config` - aktualizuje konfigurację
- `DELETE /api/availability-config` - usuwa konfigurację
- `GET /api/availability-config/check` - sprawdza istnienie

## Zmodyfikowane pliki

### 1. `models/user.go`

Dodano pole:

```go
AvailabilityConfig *UserAvailabilityConfig `gorm:"foreignKey:UserID" json:"availabilityConfig,omitempty"`
```

### 2. `services/match_service.go`

**Zmiany w `GenerateMatchesForOrganisation()`:**

- Dodano sprawdzanie czy użytkownik ma konfigurację dostępności
- Dodano filtrowanie par użytkowników według wspólnych dostępności
- Tylko pary ze wspólnymi dostępnościami są brane pod uwagę

**Zmiany w `TryGenerateMatchForUser()`:**

- Dodano sprawdzanie czy użytkownik ma konfigurację
- Dodano filtrowanie kandydatów według wspólnych dostępności
- Match tworzony tylko gdy istnieje wspólna dostępność

**Nowa zależność:**

```go
availConfigRepo repositories.UserAvailabilityConfigRepository
```

### 3. `main.go`

Dodano:

- Inicjalizację `UserAvailabilityConfigRepository`
- Inicjalizację `UserAvailabilityConfigService`
- Inicjalizację `UserAvailabilityConfigHandler`
- Przekazanie repository do `MatchService`
- Nowe endpointy API dla konfiguracji dostępności

### 4. `utils/response.go`

Dodano obsługę błędów:

- `"availability configuration not found"` → 404
- `"availability configuration already exists for this user"` → 409
- `"at least one availability slot must be selected"` → 400

## Nowe endpointy API

### Authenticated Users (wymagana autoryzacja)

#### `POST /api/availability-config`

Tworzy konfigurację dostępności dla zalogowanego użytkownika.

**Request body:**

```json
{
  "mondayMorning": true,
  "mondayAfternoon": false,
  "tuesdayMorning": true,
  "tuesdayAfternoon": true,
  "wednesdayMorning": false,
  "wednesdayAfternoon": true,
  "thursdayMorning": false,
  "thursdayAfternoon": false,
  "fridayMorning": true,
  "fridayAfternoon": true,
  "saturdayMorning": false,
  "saturdayAfternoon": false,
  "sundayMorning": false,
  "sundayAfternoon": false
}
```

**Response:** `201 Created`

```json
{
  "id": 1,
  "userId": 5,
  "mondayMorning": true,
  "mondayAfternoon": false,
  ...
  "createdAt": "2026-01-06T10:30:00Z",
  "updatedAt": "2026-01-06T10:30:00Z"
}
```

#### `GET /api/availability-config`

Pobiera konfigurację dostępności zalogowanego użytkownika.

**Response:** `200 OK`

#### `PUT /api/availability-config`

Aktualizuje konfigurację dostępności (partial update).

**Request body:**

```json
{
  "mondayMorning": false,
  "tuesdayAfternoon": true
}
```

**Response:** `200 OK`

#### `DELETE /api/availability-config`

Usuwa konfigurację dostępności.

**Response:** `200 OK`

```json
{
  "message": "Availability configuration deleted successfully"
}
```

#### `GET /api/availability-config/check`

Sprawdza czy użytkownik ma konfigurację.

**Response:** `200 OK`

```json
{
  "hasConfig": true
}
```

## Logika matchingu

### Algorytm GenerateMatchesForOrganisation()

1. Pobierz wszystkich potwierdzonych użytkowników organizacji
2. **Filtruj użytkowników:**
   - Pomiń adminów
   - Pomiń niepotwierdzonych
   - Pomiń użytkowników z pending match
   - **Pomiń użytkowników bez konfiguracji dostępności**
3. **Dla każdej pary użytkowników:**
   - Sprawdź czy nie byli matchowani w ostatnich 30 dniach
   - **Sprawdź czy mają wspólną dostępność (`HasCommonAvailability`)**
   - Jeśli tak, oblicz score i dodaj do listy par
4. Sortuj pary według score (najwyższy pierwszy)
5. Greedy algorithm: wybieraj najlepsze pary (każdy użytkownik tylko raz)
6. Twórz matche dla wybranych par

### Algorytm TryGenerateMatchForUser()

1. Pobierz użytkownika i zweryfikuj (potwierdzony, nie admin)
2. Sprawdź czy ma organizację i nie ma pending match
3. **Sprawdź czy ma konfigurację dostępności**
4. Pobierz kandydatów z tej samej organizacji
5. **Filtruj kandydatów:**
   - Pomiń siebie, adminów, niepotwierdzonych
   - Pomiń z pending match
   - Pomiń ostatnio matchowanych (30 dni)
   - **Pomiń bez konfiguracji dostępności**
   - **Pomiń bez wspólnej dostępności**
6. Wybierz kandydata z najwyższym score
7. Utwórz match

## Funkcje pomocnicze

### `HasCommonAvailability(config1, config2 *UserAvailabilityConfig) bool`

Sprawdza czy dwie konfiguracje mają co najmniej jeden wspólny slot.

**Przykład:**

```go
if models.HasCommonAvailability(userConfig, candidateConfig) {
    // Mogą być matchowani
}
```

### `GetCommonSlots(config1, config2 *UserAvailabilityConfig) []string`

Zwraca listę wspólnych slotów w czytelnym formacie.

**Przykład:**

```go
commonSlots := models.GetCommonSlots(config1, config2)
// ["Monday morning", "Tuesday afternoon", "Friday afternoon"]
```

### `GetAvailableSlots() []string`

Metoda na `UserAvailabilityConfig` zwracająca listę dostępnych slotów.

**Przykład:**

```go
slots := config.GetAvailableSlots()
// ["Monday morning", "Tuesday afternoon", "Wednesday afternoon", "Friday morning", "Friday afternoon"]
```

## Uruchomienie migracji

Aby zastosować nową tabelę:

```bash
# Backend automatycznie uruchomi migrację przy starcie
go run main.go
```

Migracja utworzy tabelę `user_availability_configs` z indeksem na `user_id`.

## Backward compatibility

⚠️ **Uwaga:** Stara tabela `match_availabilities` nadal istnieje w bazie danych, ale **nie jest już używana** przez nową logikę matchingu.

### Migracja istniejących użytkowników:

- Istniejący użytkownicy **muszą utworzyć** konfigurację dostępności
- Bez konfiguracji **nie będą matchowani**
- Frontend powinien wyświetlać komunikat o konieczności utworzenia konfiguracji

## Testowanie

### 1. Utworzenie konfiguracji

```bash
curl -X POST http://localhost:8080/api/availability-config \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "mondayMorning": true,
    "tuesdayAfternoon": true,
    "fridayMorning": true
  }'
```

### 2. Sprawdzenie matchingu

```bash
# Admin generuje matche
curl -X POST http://localhost:8080/api/admin/matches/generate \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### 3. Weryfikacja

- Tylko użytkownicy z konfiguracją powinni być matchowani
- Tylko pary ze wspólnymi dostępnościami powinny być utworzone

## Następne kroki (opcjonalne)

1. **Deprecation starej tabeli:** Rozważyć usunięcie `match_availabilities` w przyszłej migracji
2. **Email notifications:** Powiadomienia o konieczności utworzenia konfiguracji
3. **Analytics:** Dashboard pokazujący % użytkowników z konfiguracją
4. **Default availability:** Opcja ustawienia domyślnej dostępności dla nowych użytkowników
5. **Bulk import:** Możliwość importu konfiguracji z CSV

## Podsumowanie

✅ Nowa tabela `user_availability_configs` z migracją  
✅ Kompletny CRUD dla konfiguracji dostępności  
✅ Zaktualizowana logika matchingu z wymogiem wspólnych dostępności  
✅ Nowe endpointy API  
✅ Walidacja (min. 1 slot wymagany)  
✅ Funkcje pomocnicze do sprawdzania wspólnych dostępności  
✅ Obsługa błędów  
✅ Integracja z istniejącym systemem

**Match powstaje tylko gdy:**

1. ✅ Oboje użytkownicy mają konfigurację dostępności
2. ✅ Istnieje co najmniej jedna wspólna dostępność (dzień + pora)
