# Match Availability Feature - Przykłady użycia

## Format dostępności

**Nowy format availability:**

Dostępność jest teraz określana przez **dni tygodnia** i **pory dnia** zamiast konkretnych dat:

- **Dni tygodnia**: `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`
- **Pory dnia**: `morning` (przed południem), `afternoon` (po południu)

**Przykład:**

```json
{
  "Monday": ["morning", "afternoon"],
  "Wednesday": ["afternoon"],
  "Friday": ["morning"]
}
```

## Flow akceptacji matcha z dostępnością

### Automatyczne zwracanie availabilities

**Wszystkie endpointy zwracające Match obiekt automatycznie zawierają availabilities:**

- `GET /api/matches/current` - obecny pending match z availabilities
- `GET /api/matches/history` - historia matchów z availabilities
- `PATCH /api/matches/{id}/accept` - zaakceptowany match z availabilities
- `GET /api/admin/matches` - wszystkie matche organizacji z availabilities

Pole `availabilities` jest tablicą zawierającą 0-2 elementy:

- **0 elementów**: nikt jeszcze nie zaakceptował z availability
- **1 element**: jeden użytkownik zaakceptował
- **2 elementy**: obaj użytkownicy zaakceptowali

### 1. User1 akceptuje match i podaje swoją dostępność

```bash
curl -X PATCH http://localhost:8080/api/matches/1/accept \
  -H "Authorization: Bearer <user1_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "availability": {
      "Monday": ["morning", "afternoon"],
      "Wednesday": ["afternoon"],
      "Friday": ["morning"]
    }
  }'
```

**Odpowiedź:**

```json
{
  "message": "Match accepted successfully",
  "match": {
    "id": 1,
    "user1Id": 10,
    "user2Id": 15,
    "user1Accepted": true,
    "user2Accepted": false,
    "user1AcceptedAt": "2025-02-16T10:30:00Z",
    "user2AcceptedAt": null,
    "status": "pending",
    "availabilities": [
      {
        "id": 1,
        "matchId": 1,
        "userId": 10,
        "availability": {
          "Monday": ["morning", "afternoon"],
          "Wednesday": ["afternoon"],
          "Friday": ["morning"]
        }
      }
    ]
  }
}
```

### 2. User2 akceptuje match i podaje swoją dostępność

```bash
curl -X PATCH http://localhost:8080/api/matches/1/accept \
  -H "Authorization: Bearer <user2_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "availability": {
      "Monday": ["afternoon"],
      "Tuesday": ["morning"],
      "Friday": ["morning", "afternoon"]
    }
  }'
```

**Odpowiedź:**

```json
{
  "message": "Match accepted successfully",
  "match": {
    "id": 1,
    "user1Id": 10,
    "user2Id": 15,
    "user1Accepted": true,
    "user2Accepted": true,
    "user1AcceptedAt": "2025-02-16T10:30:00Z",
    "user2AcceptedAt": "2025-02-16T11:15:00Z",
    "status": "accepted",
    "expiresAt": "2025-02-21T11:15:00Z",
    "availabilities": [
      {
        "id": 1,
        "matchId": 1,
        "userId": 10,
        "availability": {
          "Monday": ["morning", "afternoon"],
          "Wednesday": ["afternoon"],
          "Friday": ["morning"]
        }
      },
      {
        "id": 2,
        "matchId": 1,
        "userId": 15,
        "availability": {
          "Monday": ["afternoon"],
          "Tuesday": ["morning"],
          "Friday": ["morning", "afternoon"]
        }
      }
    ]
  }
}
```

### 3. Pobieranie dostępności obu użytkowników

```bash
curl -X GET http://localhost:8080/api/matches/1/availabilities \
  -H "Authorization: Bearer <token>"
```

**Odpowiedź:**

```json
{
  "availabilities": [
    {
      "id": 1,
      "matchId": 1,
      "userId": 10,
      "user": {
        "id": 10,
        "firstName": "John",
        "lastName": "Doe",
        "email": "john@example.com"
      },
      "availability": {
        "Monday": ["morning", "afternoon"],
        "Wednesday": ["afternoon"],
        "Friday": ["morning"]
      },
      "createdAt": "2025-02-16T10:30:00Z",
      "updatedAt": "2025-02-16T10:30:00Z"
    },
    {
      "id": 2,
      "matchId": 1,
      "userId": 15,
      "user": {
        "id": 15,
        "firstName": "Jane",
        "lastName": "Smith",
        "email": "jane@example.com"
      },
      "availability": {
        "Monday": ["afternoon"],
        "Tuesday": ["morning"],
        "Friday": ["morning", "afternoon"]
      },
      "createdAt": "2025-02-16T11:15:00Z",
      "updatedAt": "2025-02-16T11:15:00Z"
    }
  ]
}
```

## Format danych

### Struktura availability

```json
{
  "availability": {
    "Monday": ["morning", "afternoon"],
    "Tuesday": ["afternoon"],
    "Friday": ["morning"]
  }
}
```

**Klucze (dni tygodnia):**

- Format: nazwa dnia po angielsku
- Dozwolone wartości: `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`
- Przykład: `"Monday"`, `"Friday"`
- Nie zależy od konkretnych dat w kalendarzu

**Wartości (pory dnia):**

- Format: string
- Dozwolone wartości: `"morning"` (przed południem), `"afternoon"` (po południu)
- Przykład: `["morning"]`, `["morning", "afternoon"]`
- Prosty wybór dla użytkownika

## Przykłady użycia na frontendzie

### React/TypeScript

```typescript
interface Availability {
  [weekday: string]: string[]; // {"Monday": ["morning", "afternoon"]}
}

interface MatchAcceptRequest {
  availability: Availability;
}

// Wysyłanie dostępności
const acceptMatch = async (matchId: number, availability: Availability) => {
  const response = await fetch(`/api/matches/${matchId}/accept`, {
    method: "PATCH",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ availability }),
  });

  return response.json();
};

// Przykładowe dane
const myAvailability: Availability = {
  Monday: ["morning", "afternoon"],
  Wednesday: ["afternoon"],
  Friday: ["morning"],
};

await acceptMatch(1, myAvailability);
```

### Wyświetlanie dostępności

```typescript
const AvailabilityDisplay = ({
  availabilities,
}: {
  availabilities: MatchAvailability[];
}) => {
  const periodLabels = {
    morning: "Przed południem",
    afternoon: "Po południu",
  };

  return (
    <div>
      {availabilities.map((avail) => (
        <div key={avail.id}>
          <h3>
            {avail.user.firstName} {avail.user.lastName}
          </h3>
          {Object.entries(avail.availability).map(([weekday, periods]) => (
            <div key={weekday}>
              <strong>{weekday}</strong>
              <ul>
                {periods.map((period) => (
                  <li key={period}>
                    {periodLabels[period as keyof typeof periodLabels] ||
                      period}
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      ))}
    </div>
  );
};
```

## Walidacja

Backend akceptuje dowolną strukturę JSONB, ale zalecane jest:

1. **Prawidłowe dni tygodnia**: Frontend powinien walidować, że klucze to: `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`
2. **Prawidłowe pory dnia**: Wartości to `"morning"` lub `"afternoon"`
3. **Minimalnie jeden dzień**: Przynajmniej jeden dzień z przynajmniej jedną porą dnia
4. **Unikalność**: Brak duplikatów w tablicy dla danego dnia

## Możliwe rozszerzenia

W przyszłości można dodać:

1. **Algorytm znajdowania wspólnych dni** - automatyczne sugestie wspólnych dni dostępności
2. **Finalizacja spotkania** - endpoint do wyboru konkretnej daty/godziny ze wspólnych dni
3. **Powiadomienia email** - gdy druga osoba doda swoją dostępność ✅ **Zaimplementowane**
4. **Edycja dostępności** - możliwość zmiany po zaakceptowaniu
5. **Dodatkowe pory dnia** - np. wieczór, rano wcześnie, późne popołudnie

## Struktury bazy danych

### Migracje

Wymagane migracje (wykonaj: `go run cmd/migrate/main.go migrate -up`):

1. **000010_create_match_availabilities_table.up.sql** - tworzy tabelę match_availabilities
2. **000011_add_acceptance_timestamps.up.sql** - dodaje kolumny timestamp do tabeli matches

### Tabela: match_availabilities

```sql
CREATE TABLE match_availabilities (
    id SERIAL PRIMARY KEY,
    match_id INTEGER NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    availability JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(match_id, user_id)
);
```

### Rozszerzenie tabeli: matches

```sql
ALTER TABLE matches
ADD COLUMN user1_accepted_at TIMESTAMP,
ADD COLUMN user2_accepted_at TIMESTAMP;
```

Dodane pola:

- `user1_accepted_at TIMESTAMP` - kiedy user1 zaakceptował (nullable)
- `user2_accepted_at TIMESTAMP` - kiedy user2 zaakceptował (nullable)
