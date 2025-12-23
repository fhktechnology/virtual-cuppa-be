# Email Templates - SendGrid Setup

This directory contains HTML templates for SendGrid email notifications used in Virtual Cuppa.

## Required Environment Variables

Add these to your `.env` file:

```env
SENDGRID_API_KEY=your_sendgrid_api_key_here
CONFIRM_CODE_TEMPLATE_ID=d-xxxxxxxxxxxxx
INVITATION_TEMPLATE_ID=d-xxxxxxxxxxxxx
MATCH_ACCEPTED_TEMPLATE_ID=d-xxxxxxxxxxxxx
```

## Templates Overview

### 1. Confirmation Code Template

**File**: `confirm-code-template.html` (if exists)  
**Template ID**: `CONFIRM_CODE_TEMPLATE_ID`  
**Purpose**: Send 6-digit confirmation code during registration/login

**Dynamic Data**:

- `{{Code}}` - 6-digit confirmation code

---

### 2. Invitation Template

**File**: `invitation-template.html` (if exists)  
**Template ID**: `INVITATION_TEMPLATE_ID`  
**Purpose**: Send invitation when admin creates a new user

**Dynamic Data**:

- `{{OrganisationName}}` - Name of the organisation

---

### 3. Match Accepted Template ⭐

**File**: `match-accepted-template.html`  
**Template ID**: `MATCH_ACCEPTED_TEMPLATE_ID`  
**Purpose**: Notify the other user when someone accepts the match

**Dynamic Data**:

- `{{MatchName}}` - Name of the user who accepted the match
- `{{Availability}}` - Formatted HTML showing the accepting user's available time slots

**Availability Format Example**:

```html
<div style="margin-bottom: 15px;">
  <strong>John Doe:</strong>
  <ul style="margin: 5px 0;">
    <li>2025-02-18: 09:30, 10:30, 11:00</li>
    <li>2025-02-19: 14:00, 15:00</li>
  </ul>
</div>
```

## Setup Instructions

### Creating a Template in SendGrid

1. **Login to SendGrid**

   - Go to https://sendgrid.com/
   - Navigate to Email API → Dynamic Templates

2. **Create New Template**

   - Click "Create a Dynamic Template"
   - Give it a descriptive name (e.g., "Match Accepted Notification")
   - Click "Create"

3. **Add Version**

   - Click "Add Version"
   - Choose "Code Editor"
   - Copy the HTML from the corresponding `.html` file
   - Paste it into the editor
   - Click "Save"

4. **Get Template ID**

   - The template ID will be shown in the format `d-xxxxxxxxxxxxx`
   - Copy this ID to your `.env` file

5. **Test Template**
   - Use the "Test Data" section to add sample dynamic data
   - Click "Send Test" to verify the template renders correctly

### Testing Dynamic Data

For **match-accepted-template.html**, use this test data in SendGrid:

```json
{
  "MatchName": "Jane Smith",
  "Availability": "<div style='margin-bottom: 15px;'><strong>John Doe:</strong><ul style='margin: 5px 0;'><li>2025-02-18: 09:30, 10:30, 11:00</li><li>2025-02-19: 14:00, 15:00</li></ul></div><div style='margin-bottom: 15px;'><strong>Jane Smith:</strong><ul style='margin: 5px 0;'><li>2025-02-18: 10:30, 11:00, 14:00</li><li>2025-02-20: 10:00, 11:30</li></ul></div>"
}
```

## Template Features

### Match Accepted Template Features:

- ✅ Responsive design (works on mobile and desktop)
- ✅ Gradient header with brand colors
- ✅ Clear call-to-action button
- ✅ Highlighted availability section
- ✅ Helpful tips for virtual meetings
- ✅ Professional footer

## Color Scheme

- **Primary Gradient**: `#667eea` to `#764ba2`
- **Background**: `#f4f4f4`
- **Card Background**: `#ffffff`
- **Text Primary**: `#333333`
- **Text Secondary**: `#666666`
- **Accent (Tips)**: `#f57c00` with `#fff8e1` background

## Support

If you need to modify the templates:

1. Edit the HTML file in this directory
2. Copy the updated HTML to SendGrid
3. Test thoroughly before using in production
4. Make sure all dynamic variables (`{{VariableName}}`) are preserved
