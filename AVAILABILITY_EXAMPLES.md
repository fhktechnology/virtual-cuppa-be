# Match Availability Feature - Przykłady użycia

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
      "2025-02-18": ["09:30", "10:30", "11:00"],
      "2025-02-19": ["14:00", "15:00"],
      "2025-02-20": ["10:00", "11:30", "16:00"]
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
          "2025-02-18": ["09:30", "10:30", "11:00"],
          "2025-02-19": ["14:00", "15:00"],
          "2025-02-20": ["10:00", "11:30", "16:00"]
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
      "2025-02-18": ["10:30", "11:00"],
      "2025-02-19": ["09:00", "14:00"],
      "2025-02-21": ["15:00", "16:00"]
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
    "availabilities": [
      {
        "id": 1,
        "matchId": 1,
        "userId": 10,
        "availability": {
          "2025-02-18": ["09:30", "10:30", "11:00"],
          "2025-02-19": ["14:00", "15:00"],
          "2025-02-20": ["10:00", "11:30", "16:00"]
        }
      },
      {
        "id": 2,
        "matchId": 1,
        "userId": 15,
        "availability": {
          "2025-02-18": ["10:30", "11:00"],
          "2025-02-19": ["09:00", "14:00"],
          "2025-02-21": ["15:00", "16:00"]
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
        "2025-02-18": ["09:30", "10:30", "11:00"],
        "2025-02-19": ["14:00", "15:00"],
        "2025-02-20": ["10:00", "11:30", "16:00"]
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
        "2025-02-18": ["10:30", "11:00"],
        "2025-02-19": ["09:00", "14:00"],
        "2025-02-21": ["15:00", "16:00"]
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
    "YYYY-MM-DD": ["HH:MM", "HH:MM", ...],
    "YYYY-MM-DD": ["HH:MM", "HH:MM", ...]
  }
}
```

**Klucze (daty):**

- Format: `YYYY-MM-DD` (ISO 8601)
- Przykład: `"2025-02-18"`
- Ułatwia parsowanie na frontendzie
- Sortowanie chronologiczne działa automatycznie

**Wartości (godziny):**

- Format: `HH:MM` (24-godzinny)
- Przykład: `"09:30"`, `"14:00"`, `"16:30"`
- Łatwe do wyświetlenia i konwersji na frontend
- Możliwość łatwego formatowania (np. `"09:30"` → `"9:30 AM"`)

## Przykłady użycia na frontendzie

### React/TypeScript

```typescript
interface Availability {
  [date: string]: string[]; // {"2025-02-18": ["09:30", "10:30"]}
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
  "2025-02-18": ["09:30", "10:30", "11:00"],
  "2025-02-19": ["14:00", "15:00"],
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
  return (
    <div>
      {availabilities.map((avail) => (
        <div key={avail.id}>
          <h3>
            {avail.user.firstName} {avail.user.lastName}
          </h3>
          {Object.entries(avail.availability).map(([date, times]) => (
            <div key={date}>
              <strong>{new Date(date).toLocaleDateString()}</strong>
              <ul>
                {times.map((time) => (
                  <li key={time}>{formatTime(time)}</li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      ))}
    </div>
  );
};

const formatTime = (time: string) => {
  const [hours, minutes] = time.split(":");
  const h = parseInt(hours);
  const period = h >= 12 ? "PM" : "AM";
  const displayHour = h > 12 ? h - 12 : h === 0 ? 12 : h;
  return `${displayHour}:${minutes} ${period}`;
};
```

## Walidacja

Backend akceptuje dowolną strukturę JSONB, ale zalecane jest:

1. **Daty w przyszłości**: Frontend powinien walidować, że daty są >= dzisiaj
2. **Format czasu**: `HH:MM` w formacie 24-godzinnym
3. **Minimalnie jedna data**: Przynajmniej jedna data z przynajmniej jednym słotem
4. **Unikalność godzin**: Brak duplikatów w tablicy godzin dla danego dnia

## Możliwe rozszerzenia

W przyszłości można dodać:

1. **Algorytm znajdowania wspólnych slotów** - automatyczne sugestie
2. **Finalizacja spotkania** - endpoint do wyboru konkretnej daty/godziny
3. **Powiadomienia email** - gdy druga osoba doda swoją dostępność
4. **Edycja dostępności** - możliwość zmiany po zaakceptowaniu
5. **Strefy czasowe** - obsługa różnych stref czasowych

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
