# Tania Core: Complete Rewrite Plan - Go/Vue.js to Phoenix/Elixir + Nerves

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Current Architecture Analysis](#2-current-architecture-analysis)
3. [Target Architecture](#3-target-architecture)
4. [Elixir Project Structure](#4-elixir-project-structure)
5. [Domain Model Mapping](#5-domain-model-mapping)
6. [Database Schema (Ecto Migrations)](#6-database-schema-ecto-migrations)
7. [Phoenix Contexts & Boundaries](#7-phoenix-contexts--boundaries)
8. [API & LiveView Routes](#8-api--liveview-routes)
9. [Authentication & Authorization](#9-authentication--authorization)
10. [Real-Time Features with LiveView](#10-real-time-features-with-liveview)
11. [Background Jobs & Scheduling](#11-background-jobs--scheduling)
12. [File Uploads](#12-file-uploads)
13. [Internationalization (i18n)](#13-internationalization-i18n)
14. [Nerves IoT / Raspberry Pi Support](#14-nerves-iot--raspberry-pi-support)
15. [VPS Deployment](#15-vps-deployment)
16. [Testing Strategy](#16-testing-strategy)
17. [Migration Phases](#17-migration-phases)
18. [Dependencies & Libraries](#18-dependencies--libraries)
19. [Configuration Management](#19-configuration-management)
20. [Open Questions & Decisions](#20-open-questions--decisions)

---

## 1. Executive Summary

**Tania** is an open-source farm management system currently built with:
- **Backend**: Go 1.19 (Echo v4 framework, DDD + CQRS + Event Sourcing)
- **Frontend**: Vue.js 2 SPA (Bootstrap-Vue, Vuex, Pug templates, Webpack 4) - fully integrated with backend
- **Database**: SQLite (default) / MySQL (optional)
- **License**: Apache 2.0
- **Source branch**: `master` (the complete, working v1.7.2 application)

**Rewrite targets:**
- **Phoenix Framework** (latest) with **LiveView** replacing both the Go backend and Vue.js frontend
- **Ecto** with **PostgreSQL** (primary) and **SQLite** (Nerves/embedded)
- **Nerves Project** for Raspberry Pi deployment
- **Standard VPS deployment** via Docker or bare releases

**Why Phoenix/Elixir is a good fit:**
- The existing Go codebase already uses Event Sourcing + CQRS, which maps naturally to Elixir's process model and GenServer patterns
- Phoenix LiveView eliminates the need for a separate Vue.js SPA + Webpack build pipeline, reducing complexity dramatically
- Nerves enables true embedded deployment on Raspberry Pi for on-farm use
- OTP supervision trees provide fault tolerance critical for agricultural environments
- PubSub provides real-time features that the current REST-only app lacks
- SQLite support via Ecto works on both VPS and Nerves targets
- Gettext is built into Phoenix and uses the same .po file format as the existing vue-gettext translations

---

## 2. Current Architecture Analysis

### 2.1 Backend (Go)

**Architecture Pattern**: Domain-Driven Design + CQRS + Event Sourcing

**Domain Modules** (4 bounded contexts):
| Domain | Purpose | Key Entities |
|--------|---------|--------------|
| **Assets** | Farm infrastructure | Farm, Area, Reservoir, Material |
| **Growth** | Crop lifecycle | Crop (Batch), CropActivity |
| **Tasks** | Operations management | Task |
| **User** | Authentication | User, UserAuth |

**Code Statistics**: ~225 Go files, ~32,000 lines

**Key Technical Details:**
- Echo v4 HTTP framework with middleware (CORS, auth, caching)
- EventBus (asaskevich) for async event dispatch
- Dual persistence: Event tables (write) + Read tables (denormalized queries)
- Three storage backends: SQLite, MySQL, In-Memory
- bcrypt password hashing, OAuth2 implicit grant flow
- File uploads stored on local filesystem
- Go channels for async repository operations
- Default user auto-created on startup (tania/tania)
- Demo mode (skips auth)

### 2.2 Frontend (Vue.js 2 - `master` branch)

**Status**: Fully working SPA, fully integrated with the Go backend API.

**Tech Stack:**
- Vue.js 2 with Vue Router (hash mode) and Vuex (persisted to localStorage)
- Bootstrap 4 + Bootstrap-Vue for UI components
- Pug (Jade) templates throughout all components
- Webpack 4 (raw config, no vue-cli) with code splitting, HMR
- VeeValidate v2 for form validation
- Leaflet (vue2-leaflet) for map integration
- vue-gettext for i18n
- Axios for HTTP with Bearer token auth interceptor
- moment-timezone for date formatting
- vue-toasted for notifications
- vue-slider-component for quantity sliders
- vuejs-datepicker for date pickers
- SWPrecacheWebpackPlugin for offline/PWA capability

**Frontend lives in:** `resources/js/` (built output goes to `public/`, served as static files by Go)

**Pages (16 routes):**
| Path | Component | Description |
|------|-----------|-------------|
| `/` | `home.vue` | Dashboard with stats, active crops, recent tasks |
| `/auth/login` | `auth/login.vue` | OAuth2 login |
| `/intro/farm` | `intro/farm.vue` | Onboarding step 1: create first farm (with map) |
| `/intro/reservoir` | `intro/reservoir.vue` | Onboarding step 2: create first reservoir |
| `/intro/area` | `intro/area.vue` | Onboarding step 3: create first area |
| `/areas` | `farms/areas.vue` | Area listing (card grid) |
| `/areas/:id` | `farms/area.vue` | Area detail + crops + notes + tasks |
| `/reservoirs` | `farms/reservoirs.vue` | Reservoir listing (table) |
| `/reservoirs/:id` | `farms/reservoir.vue` | Reservoir detail + notes + tasks |
| `/crops` | `farms/crops.vue` | Crop batches (Active/Archives tabs, paginated) |
| `/crop/:id` | `farms/crop.vue` | Crop detail + activity timeline + actions |
| `/crop/notes/:id` | `farms/crop-notes.vue` | Crop tasks & notes tab |
| `/tasks` | `tasks/task.vue` | Central task management with filters |
| `/materials` | `inventories/materials.vue` | Materials/inventory listing |
| `/settings/account` | `settings/account.vue` | Change password |

**Vuex Store Modules (6 modules, 41 mutations):**
- `user` - Auth state, login/logout, intro completion
- `intro` - Onboarding wizard state (farm -> reservoir -> area)
- `locations` - Countries/cities for farm location
- `farm` (aggregates 6 submodules):
  - `farms/farm` - Farm CRUD, types, current farm selection
  - `farms/area` - Area CRUD, notes, area crops
  - `farms/reservoir` - Reservoir CRUD, notes
  - `farms/crop` - Crop CRUD, notes, move/harvest/dump/water/photo, activities
  - `farms/task` - Task CRUD, filtering, complete/due status
  - `farms/inventories` - Materials CRUD, plant types, agrochemicals

**Key UI Patterns:**
- 3-step onboarding wizard forces new users to create farm -> reservoir -> area before accessing main app
- Modal-based CRUD for most entities
- Crop actions via modals: Move (quantity slider), Harvest (partial/all + weight), Dump (quantity slider), Photo upload
- Watering modal: select all or specific crops in an area
- Task filtering sidebar: by category, priority, status (completed/incomplete/overdue/today/this week/this month)
- 8 material sub-forms for different inventory types (seed, agrochemical, growing medium, etc.)
- Activity timeline on crop detail showing all lifecycle events

**i18n**: 4 languages via vue-gettext: English (default), Indonesian, Hungarian, Brazilian Portuguese
**Styling**: SCSS with Bootstrap 4 source imported, custom variables, FontAwesome 5, Material Icons

### 2.3 Database Schema (SQLite)

Uses event sourcing with paired tables:
- `*_EVENT` tables: Store versioned JSON/BLOB events per aggregate
- `*_READ` tables: Denormalized read models updated by event subscribers

**Tables**: FARM, AREA, RESERVOIR, MATERIAL, CROP (with 6 sub-tables), TASK, USER, USER_AUTH

### 2.4 API Endpoints (42 total, from Vue.js API service layer)

| Group | Endpoints | Description |
|-------|-----------|-------------|
| Auth | 2 | `POST /api/authorize` (OAuth2 login), `POST /api/user/change_password` |
| Farms | 3 | List, Create, Get types |
| Areas | 6 | Create (multipart), List, Get, Update (multipart), Create/Delete notes |
| Reservoirs | 6 | Create, List, Get, Update, Create/Delete notes |
| Materials | 5 | List (paginated), Get plant types, Get agrochemicals, Create, Update |
| Crops | 12 | Create, List (by area/farm), Get, Update, Move, Harvest, Dump, Water, Photo upload, Notes CRUD, Activities, Information/stats |
| Tasks | 6 | Create, Update, List (paginated), Search (by domain/asset or filters), Mark due, Mark complete |
| Locations | 2 | Countries, Cities by country |

### 2.5 Onboarding Flow (unique feature to preserve)

New users are forced through a 3-step wizard before accessing the main app:
1. **Create Farm** - name, description, type, map location (Leaflet + OpenStreetMap), country, city
2. **Create Reservoir** - name, source type (TAP/BUCKET), capacity
3. **Create Area** - name, size, type, location, reservoir, photo

The Vue Router guard enforces strict step ordering. The `IsNewUser` flag (`haveFarms === false`) controls this. After completing all 3 steps, `USER_COMPLETED_INTRO` is dispatched and the user gains full app access.

---

## 3. Target Architecture

```
+------------------------------------------------------------------+
|                      Phoenix Application                          |
|                                                                   |
|  +------------------+  +------------------+  +-----------------+  |
|  | Phoenix LiveView |  | REST API (JSON)  |  | Phoenix PubSub  |  |
|  | (Full UI)        |  | (Mobile/3rd party)|  | (Real-time)     |  |
|  +--------+---------+  +--------+---------+  +--------+--------+  |
|           |                     |                      |          |
|  +--------v---------------------v----------------------v--------+ |
|  |                    Phoenix Contexts                           | |
|  |  +----------+ +----------+ +---------+ +------+ +----------+ | |
|  |  | Farming  | | Growth   | | Tasks   | | Auth | | Settings | | |
|  |  +----------+ +----------+ +---------+ +------+ +----------+ | |
|  +--------------------------------------------------------------+ |
|           |                                                       |
|  +--------v-----------------------------------------------------+ |
|  |                    Ecto + PostgreSQL/SQLite                    | |
|  +--------------------------------------------------------------+ |
|           |                                                       |
|  +--------v-----------------------------------------------------+ |
|  |         Oban (Background Jobs & Scheduling)                   | |
|  +--------------------------------------------------------------+ |
+------------------------------------------------------------------+

Deployment Targets:
  [VPS / Docker]  <-->  [Raspberry Pi / Nerves]
   PostgreSQL              SQLite (embedded)
   Standard release        Firmware image
```

### 3.1 Key Architecture Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Frontend | Phoenix LiveView | Eliminates Vue.js SPA + Vuex + Webpack pipeline; server-rendered with real-time updates |
| CSS Framework | Tailwind CSS | Modern, utility-first; better fit for LiveView than Bootstrap-Vue |
| Database (VPS) | PostgreSQL | Best Ecto support, robust for production |
| Database (Nerves) | SQLite3 | Embedded, no external process needed |
| Background Jobs | Oban | Reliable job processing, cron scheduling, works with Ecto |
| Auth | phx_gen_auth + bcrypt | Phoenix's built-in auth generator |
| File Uploads | Local filesystem + optional S3 | Local for Nerves, S3 option for VPS |
| Event Sourcing | **Drop for v1** | Simplify to standard CRUD with Ecto; add event log table for audit trail |
| CQRS | **Drop for v1** | Standard Ecto queries are sufficient; no separate read models |
| i18n | Gettext | Built into Phoenix, supports .po files (existing translations reusable!) |
| API | JSON API alongside LiveView | For mobile apps / third-party integrations |

**Why drop Event Sourcing/CQRS for the rewrite:**
The original ES/CQRS adds significant complexity. For a farm management tool, standard CRUD with an audit log table provides the same practical benefits (history tracking) with far less code. This can be revisited later if needed.

---

## 4. Elixir Project Structure

```
tania/
├── mix.exs                          # Project definition & dependencies
├── config/
│   ├── config.exs                   # Shared configuration
│   ├── dev.exs                      # Development config
│   ├── test.exs                     # Test config
│   ├── prod.exs                     # Production config
│   ├── runtime.exs                  # Runtime config (env vars)
│   └── target.exs                   # Nerves target config (optional)
├── lib/
│   ├── tania/                       # Business logic (contexts)
│   │   ├── application.ex           # OTP Application
│   │   ├── repo.ex                  # Ecto Repo
│   │   ├── mailer.ex                # Email (future)
│   │   │
│   │   ├── accounts/                # Auth context
│   │   │   ├── user.ex              # User schema
│   │   │   ├── user_token.ex        # Session tokens
│   │   │   └── accounts.ex          # Context module
│   │   │
│   │   ├── farming/                 # Farm & infrastructure context
│   │   │   ├── farm.ex              # Farm schema
│   │   │   ├── area.ex              # Area schema
│   │   │   ├── reservoir.ex         # Reservoir schema
│   │   │   ├── note.ex              # Polymorphic note schema
│   │   │   └── farming.ex           # Context module
│   │   │
│   │   ├── inventory/               # Materials/inventory context
│   │   │   ├── material.ex          # Material schema
│   │   │   └── inventory.ex         # Context module
│   │   │
│   │   ├── growth/                  # Crop lifecycle context
│   │   │   ├── crop.ex              # Crop batch schema
│   │   │   ├── crop_movement.ex     # Area movement tracking
│   │   │   ├── crop_harvest.ex      # Harvest records
│   │   │   ├── crop_dump.ex         # Discard records
│   │   │   ├── crop_activity.ex     # Activity log
│   │   │   ├── crop_photo.ex        # Photo schema
│   │   │   └── growth.ex            # Context module
│   │   │
│   │   ├── tasks/                   # Task management context
│   │   │   ├── task.ex              # Task schema
│   │   │   └── tasks.ex             # Context module
│   │   │
│   │   ├── locations/               # Geographic data context
│   │   │   └── locations.ex         # Countries/cities helpers
│   │   │
│   │   └── audit/                   # Audit trail (replaces event sourcing)
│   │       ├── log_entry.ex         # Audit log schema
│   │       └── audit.ex             # Context module
│   │
│   ├── tania_web/                   # Web layer
│   │   ├── endpoint.ex              # Phoenix Endpoint
│   │   ├── router.ex                # All routes
│   │   ├── gettext.ex               # i18n
│   │   ├── telemetry.ex             # Metrics
│   │   │
│   │   ├── components/              # Shared UI components
│   │   │   ├── core_components.ex   # Buttons, modals, tables, forms
│   │   │   ├── layouts.ex           # App layout with sidebar
│   │   │   └── icons.ex             # SVG icon components
│   │   │
│   │   ├── live/                    # LiveView pages
│   │   │   ├── dashboard_live.ex
│   │   │   ├── intro_live/          # Onboarding wizard (3 steps)
│   │   │   │   ├── farm.ex          # Step 1: Create farm (with map)
│   │   │   │   ├── reservoir.ex     # Step 2: Create reservoir
│   │   │   │   └── area.ex          # Step 3: Create area
│   │   │   ├── farm_live/
│   │   │   │   ├── index.ex         # Farm list
│   │   │   │   ├── show.ex          # Farm detail
│   │   │   │   └── form_component.ex
│   │   │   ├── area_live/
│   │   │   │   ├── index.ex
│   │   │   │   ├── show.ex
│   │   │   │   └── form_component.ex
│   │   │   ├── reservoir_live/
│   │   │   │   ├── index.ex
│   │   │   │   ├── show.ex
│   │   │   │   └── form_component.ex
│   │   │   ├── material_live/
│   │   │   │   ├── index.ex
│   │   │   │   ├── show.ex
│   │   │   │   └── form_component.ex
│   │   │   ├── crop_live/
│   │   │   │   ├── index.ex           # Active/Archives tabs, paginated
│   │   │   │   ├── show.ex            # Detail + activity timeline
│   │   │   │   ├── notes.ex           # Tasks & notes tab
│   │   │   │   ├── form_component.ex  # Create/edit batch form
│   │   │   │   ├── harvest_component.ex # Partial/all + weight
│   │   │   │   ├── move_component.ex    # Area transfer + quantity slider
│   │   │   │   ├── dump_component.ex    # Discard + quantity slider
│   │   │   │   └── upload_component.ex  # Photo upload + description
│   │   │   ├── task_live/
│   │   │   │   ├── index.ex           # Central tasks + filter sidebar
│   │   │   │   ├── form_component.ex  # General task form
│   │   │   │   ├── crop_task_form_component.ex # Crop-specific (nutrient/pest)
│   │   │   │   ├── water_task_component.ex     # Watering modal (all/partial)
│   │   │   │   └── task_list_component.ex      # Reusable task list (used in area, reservoir, crop detail)
│   │   │   └── account_live/
│   │   │       ├── settings.ex
│   │   │       ├── login.ex
│   │   │       └── registration.ex
│   │   │
│   │   └── controllers/             # JSON API controllers
│   │       ├── api/
│   │       │   ├── farm_controller.ex
│   │       │   ├── area_controller.ex
│   │       │   ├── crop_controller.ex
│   │       │   ├── task_controller.ex
│   │       │   └── ...
│   │       ├── user_session_controller.ex
│   │       └── fallback_controller.ex
│   │
│   └── tania_web.ex                 # Web module macros
│
├── priv/
│   ├── repo/migrations/             # Ecto migrations
│   ├── static/                      # Static assets (CSS, JS, images)
│   ├── gettext/                     # Translation .po files
│   └── uploads/                     # Local file uploads
│
├── assets/
│   ├── css/app.css                  # Tailwind CSS
│   ├── js/app.js                    # Phoenix JS hooks
│   ├── js/hooks/                    # LiveView JS hooks
│   │   ├── leaflet_map.js           # Leaflet map for farm location (replaces vue2-leaflet)
│   │   ├── quantity_slider.js       # Range slider for crop move/dump (replaces vue-slider-component)
│   │   └── datepicker.js            # Date picker for tasks (or use native HTML date input)
│   └── tailwind.config.js
│
├── test/
│   ├── tania/                       # Context tests
│   ├── tania_web/                   # Web layer tests
│   └── support/                     # Test helpers & fixtures
│
└── nerves/                          # Nerves-specific (see Section 14)
    ├── config/target.exs
    └── rootfs_overlay/
```

---

## 5. Domain Model Mapping

### 5.1 Go to Elixir Entity Mapping

#### Farm

| Go Field | Elixir Field | Type | Notes |
|----------|-------------|------|-------|
| UID (uuid) | id (binary_id) | Ecto.UUID | Primary key |
| Name | name | :string | Required, validated |
| Type | type | Ecto.Enum | :organic, :hydroponic, :aquaponic, :mushroom, :livestock, :fisheries, :permaculture |
| Latitude | latitude | :float | |
| Longitude | longitude | :float | |
| Country | country | :string | |
| City | city | :string | |
| IsActive | is_active | :boolean | Default: true |
| CreatedDate | inserted_at | :utc_datetime | Ecto timestamp |
| - | updated_at | :utc_datetime | Ecto timestamp |

**Associations**: has_many :areas, has_many :reservoirs, has_many :crops

#### Area

| Go Field | Elixir Field | Type | Notes |
|----------|-------------|------|-------|
| UID | id | Ecto.UUID | Primary key |
| Name | name | :string | Required |
| Size | size | :float | |
| SizeUnit | size_unit | Ecto.Enum | :sqm, :hectare |
| Type | type | Ecto.Enum | :seeding, :growing |
| Location | location | Ecto.Enum | :indoor, :outdoor |
| Photo* | photo | :string | File path/URL |
| FarmUID | farm_id | Ecto.UUID | belongs_to :farm |
| ReservoirUID | reservoir_id | Ecto.UUID | belongs_to :reservoir (optional) |

**Associations**: belongs_to :farm, belongs_to :reservoir, has_many :notes (polymorphic), has_many :crops

#### Reservoir

| Go Field | Elixir Field | Type | Notes |
|----------|-------------|------|-------|
| UID | id | Ecto.UUID | Primary key |
| Name | name | :string | Required |
| WaterSourceType | water_source_type | Ecto.Enum | :tap, :well, :rain_water, :river, :tank |
| WaterSourceCapacity | water_source_capacity | :float | Liters |
| FarmUID | farm_id | Ecto.UUID | belongs_to :farm |

**Associations**: belongs_to :farm, has_many :areas, has_many :notes (polymorphic)

#### Material (Inventory)

| Go Field | Elixir Field | Type | Notes |
|----------|-------------|------|-------|
| UID | id | Ecto.UUID | Primary key |
| Name | name | :string | Required |
| PricePerUnit | price_per_unit | :decimal | |
| CurrencyCode | currency_code | :string | Default: "EUR" |
| Type | type | Ecto.Enum | :seed, :agrochemical, :growing_medium, :label_crop_support, :seeding_container, :post_harvest, :tool, :other |
| TypeData | type_data | :map | Embedded: plant_type for seeds |
| Quantity | quantity | :float | |
| QuantityUnit | quantity_unit | Ecto.Enum | :seeds, :packets, :gram, :kilogram, :bags, :bottles, :cubic_metre, :pieces, :units |
| ExpirationDate | expiration_date | :date | |
| Notes | notes | :string | |
| ProducedBy | produced_by | :string | |

#### Crop (Batch)

| Go Field | Elixir Field | Type | Notes |
|----------|-------------|------|-------|
| UID | id | Ecto.UUID | Primary key |
| BatchID | batch_id | :string | Auto-generated |
| Status | status | Ecto.Enum | :active, :archived |
| Type | type | Ecto.Enum | :seeding, :growing |
| ContainerQuantity | container_quantity | :integer | |
| ContainerType | container_type | Ecto.Enum | :tray, :pot |
| ContainerCell | container_cell | :integer | Cells per tray |
| InventoryUID | material_id | Ecto.UUID | belongs_to :material |
| FarmUID | farm_id | Ecto.UUID | belongs_to :farm |
| InitialArea* | initial_area_id | Ecto.UUID | belongs_to :area |
| AreaStatusSeeding | seeding_quantity | :integer | Computed |
| AreaStatusGrowing | growing_quantity | :integer | Computed |
| AreaStatusDumped | dumped_quantity | :integer | Computed |
| LastWatered | last_watered_at | :utc_datetime | |
| LastFertilized | last_fertilized_at | :utc_datetime | |
| LastPesticided | last_pesticided_at | :utc_datetime | |
| LastPruned | last_pruned_at | :utc_datetime | |

**Associations**: belongs_to :farm, belongs_to :material, belongs_to :initial_area (Area), has_many :movements, has_many :harvests, has_many :dumps, has_many :activities, has_many :photos, has_many :notes (polymorphic)

#### CropMovement (new name for MovedArea)

| Field | Type | Notes |
|-------|------|-------|
| id | Ecto.UUID | Primary key |
| crop_id | Ecto.UUID | belongs_to :crop |
| area_id | Ecto.UUID | belongs_to :area |
| initial_quantity | :integer | |
| current_quantity | :integer | |
| last_watered_at | :utc_datetime | |
| last_fertilized_at | :utc_datetime | |
| last_pesticided_at | :utc_datetime | |
| last_pruned_at | :utc_datetime | |

#### CropHarvest

| Field | Type | Notes |
|-------|------|-------|
| id | Ecto.UUID | Primary key |
| crop_id | Ecto.UUID | belongs_to :crop |
| source_area_id | Ecto.UUID | belongs_to :area |
| quantity | :integer | Plants harvested |
| produced_gram_quantity | :float | Weight in grams |

#### CropDump

| Field | Type | Notes |
|-------|------|-------|
| id | Ecto.UUID | Primary key |
| crop_id | Ecto.UUID | belongs_to :crop |
| source_area_id | Ecto.UUID | belongs_to :area |
| quantity | :integer | Plants discarded |

#### CropActivity

| Field | Type | Notes |
|-------|------|-------|
| id | Ecto.UUID | Primary key |
| crop_id | Ecto.UUID | belongs_to :crop |
| activity_type | Ecto.Enum | :seed, :move, :harvest, :dump, :water, :fertilize, :prune, :pesticide, :photo |
| description | :string | |

#### Task

| Go Field | Elixir Field | Type | Notes |
|----------|-------------|------|-------|
| UID | id | Ecto.UUID | Primary key |
| Title | title | :string | Required |
| Description | description | :string | |
| DueDate | due_date | :date | |
| CompletedDate | completed_at | :utc_datetime | |
| CancelledDate | cancelled_at | :utc_datetime | |
| Priority | priority | Ecto.Enum | :high, :medium, :low |
| Status | status | Ecto.Enum | :pending, :completed, :cancelled |
| Category | category | Ecto.Enum | :farm_maintenance, :inventory, :planting, :harvesting, :water_management, :pest_control, :safety, :sanitation, :general |
| Domain | domain | Ecto.Enum | :area, :crop, :material, :reservoir |
| AssetID | asset_id | Ecto.UUID | Polymorphic reference |
| IsDue | is_due | :boolean | Computed |

#### Note (Polymorphic)

| Field | Type | Notes |
|-------|------|-------|
| id | Ecto.UUID | Primary key |
| content | :string | Required |
| notable_type | :string | "Area", "Reservoir", "Crop" |
| notable_id | Ecto.UUID | Polymorphic FK |

#### User

| Go Field | Elixir Field | Type | Notes |
|----------|-------------|------|-------|
| UID | id | Ecto.UUID | Primary key |
| Username | email | :string | Switch to email-based auth |
| Password | hashed_password | :string | bcrypt via Bcrypt.hash_pwd_salt |

#### AuditLog (replaces Event Sourcing)

| Field | Type | Notes |
|-------|------|-------|
| id | :id | Auto-increment |
| user_id | Ecto.UUID | Who performed action |
| action | :string | "created", "updated", "deleted" |
| resource_type | :string | "Farm", "Crop", etc. |
| resource_id | Ecto.UUID | |
| changes | :map | JSON diff of changes |
| inserted_at | :utc_datetime | When |

---

## 6. Database Schema (Ecto Migrations)

### Migration Order

```
001_create_users.exs
002_create_farms.exs
003_create_reservoirs.exs
004_create_areas.exs
005_create_materials.exs
006_create_crops.exs
007_create_crop_movements.exs
008_create_crop_harvests.exs
009_create_crop_dumps.exs
010_create_crop_activities.exs
011_create_crop_photos.exs
012_create_notes.exs
013_create_tasks.exs
014_create_audit_logs.exs
015_create_user_tokens.exs
```

### Key Migration: Crops (most complex)

```elixir
# Example: 006_create_crops.exs
def change do
  create table(:crops, primary_key: false) do
    add :id, :binary_id, primary_key: true
    add :batch_id, :string, null: false
    add :status, :string, null: false, default: "active"
    add :type, :string, null: false
    add :container_quantity, :integer
    add :container_type, :string
    add :container_cell, :integer
    add :seeding_quantity, :integer, default: 0
    add :growing_quantity, :integer, default: 0
    add :dumped_quantity, :integer, default: 0
    add :last_watered_at, :utc_datetime
    add :last_fertilized_at, :utc_datetime
    add :last_pesticided_at, :utc_datetime
    add :last_pruned_at, :utc_datetime
    add :farm_id, references(:farms, type: :binary_id), null: false
    add :material_id, references(:materials, type: :binary_id)
    add :initial_area_id, references(:areas, type: :binary_id)
    timestamps(type: :utc_datetime)
  end

  create index(:crops, [:farm_id])
  create index(:crops, [:status])
  create index(:crops, [:material_id])
end
```

---

## 7. Phoenix Contexts & Boundaries

### 7.1 Context Responsibilities

```
Tania.Accounts
  - User registration, login, session management
  - Password hashing & validation
  - User token management
  - Generated by: mix phx.gen.auth

Tania.Farming
  - Farm CRUD (create, list, update, get)
  - Area CRUD with farm association
  - Reservoir CRUD with farm association
  - Notes management (polymorphic for areas & reservoirs)
  - Farm type/location validation

Tania.Inventory
  - Material CRUD
  - Quantity tracking
  - Expiration date management
  - Plant type categorization
  - Price tracking

Tania.Growth
  - Crop batch lifecycle (create -> grow -> harvest/dump -> archive)
  - Crop movements between areas
  - Harvest recording
  - Dump/discard recording
  - Care activity logging (water, fertilize, prune, pesticide)
  - Photo management
  - Batch ID generation
  - Days-since-seeding calculation
  - Crop notes

Tania.Tasks
  - Task CRUD
  - Status lifecycle (pending -> completed/cancelled)
  - Priority management
  - Due date tracking
  - Domain/asset linking (polymorphic)
  - Filtering (by status, priority, category, due date)

Tania.Locations
  - Country listing
  - City lookup by country
  - (Use `countries` hex package or embed data)

Tania.Audit
  - Automatic audit logging via Ecto hooks or explicit calls
  - Query audit history per resource
```

### 7.2 Context Interaction Rules

```
Farming --> (no deps)
Inventory --> (no deps)
Growth --> Farming (reads areas, farms), Inventory (reads materials)
Tasks --> Farming, Growth, Inventory (reads asset info for display)
Accounts --> (no deps)
Audit --> All contexts (receives change notifications)
```

### 7.3 Example Context Function Signatures

```elixir
# Tania.Growth context
defmodule Tania.Growth do
  def list_crops(farm_id, opts \\ [])
  def get_crop!(id)
  def create_crop(attrs)
  def update_crop(crop, attrs)
  def archive_crop(crop)
  def move_crop(crop, %{area_id: area_id, quantity: qty})
  def harvest_crop(crop, %{area_id: area_id, quantity: qty, weight_grams: w})
  def dump_crop(crop, %{area_id: area_id, quantity: qty})
  def water_crop(crop)
  def fertilize_crop(crop)
  def prune_crop(crop)
  def pesticide_crop(crop)
  def add_crop_photo(crop, photo_params)
  def add_crop_note(crop, content)
  def remove_crop_note(note_id)
  def list_crop_activities(crop_id)
  def days_since_seeding(crop)
end
```

---

## 8. API & LiveView Routes

### 8.1 LiveView Routes (Primary UI)

```elixir
# router.ex
scope "/", TaniaWeb do
  pipe_through [:browser, :require_authenticated_user]

  # Onboarding wizard (authenticated but no farm yet)
  live_session :onboarding, on_mount: [{TaniaWeb.UserAuth, :ensure_authenticated}] do
    live "/intro/farm", IntroLive.Farm, :new
    live "/intro/reservoir", IntroLive.Reservoir, :new
    live "/intro/area", IntroLive.Area, :new
  end

  # Main app (authenticated + has farms)
  live_session :authenticated, on_mount: [
    {TaniaWeb.UserAuth, :ensure_authenticated},
    {TaniaWeb.UserAuth, :ensure_has_farm}
  ] do
    live "/", DashboardLive, :index

    # Farms
    live "/farms", FarmLive.Index, :index
    live "/farms/new", FarmLive.Index, :new
    live "/farms/:id", FarmLive.Show, :show
    live "/farms/:id/edit", FarmLive.Show, :edit

    # Areas (card grid listing)
    live "/areas", AreaLive.Index, :index
    live "/areas/new", AreaLive.Index, :new
    live "/areas/:id", AreaLive.Show, :show
    live "/areas/:id/edit", AreaLive.Show, :edit

    # Reservoirs
    live "/reservoirs", ReservoirLive.Index, :index
    live "/reservoirs/new", ReservoirLive.Index, :new
    live "/reservoirs/:id", ReservoirLive.Show, :show
    live "/reservoirs/:id/edit", ReservoirLive.Show, :edit

    # Materials (8 sub-type forms)
    live "/materials", MaterialLive.Index, :index
    live "/materials/new", MaterialLive.Index, :new
    live "/materials/:id", MaterialLive.Show, :show
    live "/materials/:id/edit", MaterialLive.Show, :edit

    # Crops (Active/Archives tabs)
    live "/crops", CropLive.Index, :index
    live "/crops/new", CropLive.Index, :new
    live "/crops/:id", CropLive.Show, :show
    live "/crops/:id/harvest", CropLive.Show, :harvest
    live "/crops/:id/move", CropLive.Show, :move
    live "/crops/:id/dump", CropLive.Show, :dump
    live "/crops/:id/notes", CropLive.Notes, :index

    # Tasks (central + context-specific)
    live "/tasks", TaskLive.Index, :index
    live "/tasks/new", TaskLive.Index, :new
    live "/tasks/:id", TaskLive.Index, :show

    # Account
    live "/account/settings", AccountLive.Settings, :edit
  end
end

# Auth routes (no auth required)
scope "/", TaniaWeb do
  pipe_through [:browser, :redirect_if_user_is_authenticated]

  live_session :unauthenticated, on_mount: [{TaniaWeb.UserAuth, :redirect_if_user_is_authenticated}] do
    live "/users/log_in", AccountLive.Login, :new
    live "/users/register", AccountLive.Registration, :new
  end
end
```

**Onboarding Flow Logic:**
The `ensure_has_farm` on_mount hook checks if the user has any farms. If not, it redirects to `/intro/farm`. The intro LiveViews enforce step ordering (farm -> reservoir -> area) and redirect to `/` upon completion. This preserves the Vue.js onboarding wizard behavior.

### 8.2 JSON API Routes (for mobile/third-party)

```elixir
scope "/api", TaniaWeb.API do
  pipe_through [:api, :api_auth]

  resources "/farms", FarmController, only: [:index, :show, :create, :update] do
    resources "/areas", AreaController, only: [:index, :create]
    resources "/reservoirs", ReservoirController, only: [:index, :create]
    resources "/crops", CropController, only: [:index]
  end

  resources "/areas", AreaController, only: [:show, :update]
  resources "/reservoirs", ReservoirController, only: [:show, :update]
  resources "/materials", MaterialController, except: [:delete]

  resources "/crops", CropController, only: [:show, :create, :update] do
    post "/move", CropController, :move
    post "/harvest", CropController, :harvest
    post "/dump", CropController, :dump
    post "/water", CropController, :water
    post "/photos", CropController, :upload_photo
    get "/activities", CropController, :activities
  end

  resources "/tasks", TaskController do
    post "/complete", TaskController, :complete
    post "/cancel", TaskController, :cancel
  end

  get "/locations/countries", LocationController, :countries
  get "/locations/cities/:country", LocationController, :cities
end
```

---

## 9. Authentication & Authorization

### 9.1 Implementation

Use `mix phx.gen.auth` which generates:
- User schema with hashed_password (bcrypt)
- Session-based authentication for LiveView
- Token-based authentication for API
- Login, registration, password reset flows
- Plugs: `require_authenticated_user`, `redirect_if_user_is_authenticated`

### 9.2 API Authentication

For the JSON API, use Bearer token authentication:

```elixir
# API auth plug
defmodule TaniaWeb.API.AuthPlug do
  def call(conn, _opts) do
    with ["Bearer " <> token] <- get_req_header(conn, "authorization"),
         {:ok, user} <- Tania.Accounts.get_user_by_api_token(token) do
      assign(conn, :current_user, user)
    else
      _ -> conn |> send_resp(401, "Unauthorized") |> halt()
    end
  end
end
```

### 9.3 Demo Mode

Support a demo mode via configuration that auto-authenticates:

```elixir
# config/dev.exs
config :tania, :demo_mode, true

# In auth plug
if Application.get_env(:tania, :demo_mode, false) do
  # Auto-assign demo user
end
```

### 9.4 Default User Seeding

```elixir
# priv/repo/seeds.exs
if Tania.Accounts.list_users() == [] do
  Tania.Accounts.register_user(%{
    email: "admin@tania.local",
    password: "tania_admin_2024"
  })
end
```

---

## 10. Real-Time Features with LiveView

The current Go app has NO real-time features. Phoenix LiveView provides them for free.

### 10.1 Real-Time Updates

```elixir
# When a crop is watered, broadcast to all viewers
def water_crop(crop) do
  {:ok, updated_crop} = # ... update logic
  Phoenix.PubSub.broadcast(Tania.PubSub, "crop:#{crop.id}", {:crop_updated, updated_crop})
  Phoenix.PubSub.broadcast(Tania.PubSub, "farm:#{crop.farm_id}", {:crop_watered, updated_crop})
  {:ok, updated_crop}
end

# In CropLive.Show
def mount(%{"id" => id}, _session, socket) do
  crop = Growth.get_crop!(id)
  if connected?(socket), do: Phoenix.PubSub.subscribe(Tania.PubSub, "crop:#{id}")
  {:ok, assign(socket, crop: crop)}
end

def handle_info({:crop_updated, crop}, socket) do
  {:noreply, assign(socket, crop: crop)}
end
```

### 10.2 Dashboard Live Updates

The dashboard can show real-time stats by subscribing to farm-level PubSub topics:
- Active crops count changes
- Tasks completed
- Recent activities stream

---

## 11. Background Jobs & Scheduling

### 11.1 Oban Jobs

```elixir
# mix.exs
{:oban, "~> 2.17"}

# config/config.exs
config :tania, Oban,
  repo: Tania.Repo,
  queues: [default: 10, mailer: 5],
  plugins: [
    Oban.Plugins.Pruner,
    {Oban.Plugins.Cron, crontab: [
      {"0 6 * * *", Tania.Workers.TaskDueChecker},    # Check due tasks daily at 6am
      {"0 0 * * *", Tania.Workers.CropAgeUpdater},    # Update crop ages daily
    ]}
  ]
```

### 11.2 Worker Examples

```elixir
# Mark tasks as due when due_date passes
defmodule Tania.Workers.TaskDueChecker do
  use Oban.Worker, queue: :default

  @impl Oban.Worker
  def perform(_job) do
    Tania.Tasks.mark_overdue_tasks()
    :ok
  end
end
```

### 11.3 Nerves Consideration

Oban requires PostgreSQL. For Nerves/SQLite deployments, use `Quantum` or simple GenServer-based schedulers instead:

```elixir
# For Nerves: use Quantum scheduler
config :tania, Tania.Scheduler,
  jobs: [
    {"0 6 * * *", {Tania.Workers.TaskDueChecker, :run, []}},
  ]
```

---

## 12. File Uploads

### 12.1 Phoenix LiveView Uploads

LiveView has built-in upload support with drag-and-drop, progress bars, and validation:

```elixir
# In CropLive.Show
def mount(_params, _session, socket) do
  {:ok,
   socket
   |> allow_upload(:photo,
     accept: ~w(.jpg .jpeg .png .webp),
     max_entries: 5,
     max_file_size: 10_000_000
   )}
end

def handle_event("save_photo", _params, socket) do
  uploaded_files =
    consume_uploaded_entries(socket, :photo, fn %{path: path}, entry ->
      dest = Path.join(upload_dir(), "#{Ecto.UUID.generate()}-#{entry.client_name}")
      File.cp!(path, dest)
      {:ok, dest}
    end)

  # Save file references to database
  {:noreply, socket}
end
```

### 12.2 Storage Strategy

| Target | Storage | Config |
|--------|---------|--------|
| VPS | Local filesystem or S3 | Configurable via runtime.exs |
| Nerves | Local filesystem (SD card) | Fixed path on device |

Consider using a storage abstraction like `Waffle` (formerly Arc) for pluggable backends.

---

## 13. Internationalization (i18n)

### 13.1 Reuse Existing Translations

The current Vue.js app uses `vue-gettext` with .po files. Phoenix uses Gettext which also uses .po files. The existing translations can be migrated directly with minimal changes.

### 13.2 Supported Locales (from `master` branch)

| Locale | Language | Status |
|--------|----------|--------|
| en | English | Complete (default) |
| id | Indonesian | Complete |
| hu | Hungarian | Complete |
| pt_BR | Portuguese (Brazil) | Complete |

Note: The `2.0-dev` branch added Vietnamese, Greek, and Czech, but those are in the unfinished Next.js rewrite. The `master` branch has 4 fully working languages.

### 13.3 Implementation

```elixir
# In LiveView templates
<h1><%= gettext("Dashboard") %></h1>
<p><%= ngettext("1 crop", "%{count} crops", @crop_count) %></p>

# Locale switching via plug or LiveView hook
# Store user preference in session or user record
```

### 13.4 Migration Steps

1. Copy existing .po files from `resources/js/languages/` (master branch) to `priv/gettext/`
2. Rename to Phoenix Gettext conventions: `priv/gettext/hu/LC_MESSAGES/default.po`
3. Convert vue-gettext `<translate>` tag strings to `gettext()` function calls
4. Extract new strings with `mix gettext.extract`
5. Merge with `mix gettext.merge priv/gettext`

---

## 14. Nerves IoT / Raspberry Pi Support

### 14.1 Architecture: Poncho Project

Use a **poncho project** structure (NOT umbrella) for maximum flexibility:

```
tania-umbrella/
├── tania/              # Core business logic + Ecto (shared)
├── tania_web/          # Phoenix web interface (shared)
├── tania_vps/          # VPS-specific release config
│   ├── mix.exs         # Depends on tania + tania_web
│   └── config/
└── tania_nerves/       # Nerves firmware project
    ├── mix.exs         # Depends on tania + tania_web + nerves deps
    ├── config/
    │   ├── config.exs
    │   └── target.exs  # Hardware-specific config
    ├── rootfs_overlay/  # Custom filesystem files
    └── rel/
```

**Alternative (simpler):** Single project with conditional compilation:

```elixir
# mix.exs
defp deps do
  [
    {:phoenix, "~> 1.7"},
    {:ecto_sql, "~> 3.11"},
    {:ecto_sqlite3, "~> 0.16"},    # SQLite for all targets
    {:postgrex, "~> 0.18"},        # PostgreSQL for VPS only
  ] ++ nerves_deps()
end

defp nerves_deps do
  if Mix.target() != :host do
    [
      {:nerves, "~> 1.10"},
      {:nerves_system_rpi4, "~> 1.24"},    # Raspberry Pi 4
      {:nerves_system_rpi3, "~> 1.25"},    # Raspberry Pi 3
      {:vintage_net, "~> 0.13"},           # Networking
      {:vintage_net_wifi, "~> 0.12"},      # WiFi
      {:vintage_net_ethernet, "~> 0.11"},  # Ethernet
      {:nerves_time, "~> 0.4"},            # NTP time sync
      {:ring_logger, "~> 0.10"},           # In-memory logging
    ]
  else
    []
  end
end
```

### 14.2 Nerves-Specific Considerations

| Concern | Solution |
|---------|----------|
| **Database** | SQLite3 via `ecto_sqlite3` (no PostgreSQL on device) |
| **Storage** | SD card with wear-leveling partition for uploads |
| **Networking** | WiFi config via `VintageNet`; mDNS for discovery (`tania.local`) |
| **Time** | NTP sync via `nerves_time` (no RTC on most Pi models) |
| **Updates** | OTA firmware updates via `NervesHub` or SSH |
| **Background Jobs** | `Quantum` or GenServer (no Oban without PostgreSQL) |
| **File Uploads** | Local filesystem on data partition |
| **Logging** | `RingLogger` (in-memory circular buffer) |
| **Boot** | Auto-start Phoenix on boot, accessible via WiFi AP or LAN |

### 14.3 Nerves Configuration

```elixir
# config/target.exs (Nerves-specific)
config :tania, Tania.Repo,
  database: "/data/tania/tania.db",
  pool_size: 5

config :tania, TaniaWeb.Endpoint,
  http: [port: 80],
  url: [host: "tania.local"],
  server: true

# WiFi configuration
config :vintage_net,
  regulatory_domain: "US",
  config: [
    {"wlan0", %{type: VintageNetWiFi, ...}},
    {"eth0", %{type: VintageNetEthernet, ipv4: %{method: :dhcp}}}
  ]
```

### 14.4 Nerves Data Persistence

Create a separate data partition on the SD card for database and uploads:

```elixir
# In application startup
File.mkdir_p!("/data/tania/uploads/areas")
File.mkdir_p!("/data/tania/uploads/crops")
```

### 14.5 Hardware Integration Opportunities (Future)

Nerves opens the door for direct sensor integration:
- **GPIO**: Soil moisture sensors, temperature/humidity sensors
- **I2C/SPI**: Advanced sensors, displays
- **Camera**: Pi Camera for automated crop photos
- **Relay control**: Automated irrigation via GPIO relays

These are NOT in scope for the initial rewrite but are natural extensions.

### 14.6 Supported Raspberry Pi Models

| Model | Nerves System | RAM | Notes |
|-------|--------------|-----|-------|
| RPi 3B+ | `nerves_system_rpi3` | 1GB | Minimum viable |
| RPi 4B | `nerves_system_rpi4` | 2-8GB | Recommended |
| RPi 5 | `nerves_system_rpi5` | 4-8GB | Best performance |
| RPi Zero 2W | `nerves_system_rpi0_2` | 512MB | Ultra-compact |

---

## 15. VPS Deployment

### 15.1 Docker Deployment

```dockerfile
# Dockerfile
FROM hexpm/elixir:1.16.1-erlang-26.2.2-debian-bullseye-20240130 AS build

WORKDIR /app
ENV MIX_ENV=prod

COPY mix.exs mix.lock ./
RUN mix deps.get --only prod && mix deps.compile

COPY lib lib
COPY priv priv
COPY assets assets
COPY config config

RUN mix assets.deploy && mix compile && mix release

# --- Runtime ---
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y libstdc++6 openssl libncurses5 locales \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build /app/_build/prod/rel/tania ./

ENV DATABASE_URL=ecto://postgres:postgres@db/tania
ENV SECRET_KEY_BASE=<generate-with-mix-phx-gen-secret>
ENV PHX_HOST=your-domain.com

CMD ["bin/tania", "start"]
```

```yaml
# docker-compose.yml
services:
  db:
    image: postgres:16-alpine
    volumes:
      - pgdata:/var/lib/postgresql/data
    environment:
      POSTGRES_DB: tania
      POSTGRES_USER: tania
      POSTGRES_PASSWORD: ${DB_PASSWORD}

  app:
    build: .
    ports:
      - "80:4000"
    depends_on:
      - db
    environment:
      DATABASE_URL: ecto://tania:${DB_PASSWORD}@db/tania
      SECRET_KEY_BASE: ${SECRET_KEY_BASE}
      PHX_HOST: ${PHX_HOST}
    volumes:
      - uploads:/app/priv/uploads

volumes:
  pgdata:
  uploads:
```

### 15.2 Bare Metal Release

```bash
# Build release on matching OS/arch
MIX_ENV=prod mix deps.get
MIX_ENV=prod mix assets.deploy
MIX_ENV=prod mix release

# Deploy
scp _build/prod/rel/tania/tania-*.tar.gz server:
ssh server "tar xzf tania-*.tar.gz && bin/tania start"
```

### 15.3 systemd Service

```ini
# /etc/systemd/system/tania.service
[Unit]
Description=Tania Farm Management
After=network.target postgresql.service

[Service]
Type=exec
User=tania
WorkingDirectory=/opt/tania
ExecStart=/opt/tania/bin/tania start
ExecStop=/opt/tania/bin/tania stop
Restart=on-failure
RestartSec=5
EnvironmentFile=/opt/tania/.env

[Install]
WantedBy=multi-user.target
```

---

## 16. Testing Strategy

### 16.1 Test Structure

```
test/
├── tania/
│   ├── accounts_test.exs          # User registration, login
│   ├── farming_test.exs           # Farm, Area, Reservoir CRUD
│   ├── inventory_test.exs         # Material CRUD
│   ├── growth_test.exs            # Crop lifecycle (most complex)
│   ├── tasks_test.exs             # Task CRUD & lifecycle
│   └── audit_test.exs             # Audit logging
├── tania_web/
│   ├── live/
│   │   ├── dashboard_live_test.exs
│   │   ├── farm_live_test.exs
│   │   ├── area_live_test.exs
│   │   ├── crop_live_test.exs
│   │   ├── task_live_test.exs
│   │   └── ...
│   └── controllers/
│       └── api/
│           ├── farm_controller_test.exs
│           └── ...
└── support/
    ├── fixtures/
    │   ├── accounts_fixtures.ex
    │   ├── farming_fixtures.ex
    │   ├── growth_fixtures.ex
    │   └── ...
    ├── conn_case.ex
    └── data_case.ex
```

### 16.2 Key Test Scenarios

**Growth Context (highest complexity):**
- Create crop batch with valid/invalid params
- Move crop between areas (seeding->seeding, seeding->growing, growing->growing)
- Harvest crop (full and partial)
- Dump/discard crop
- Water/fertilize/prune/pesticide logging
- Days since seeding calculation
- Crop archival
- Activity log recording

**Task Context:**
- Create task linked to various domains (area, crop, material, reservoir)
- Complete/cancel task lifecycle
- Due date checking and overdue marking
- Filtering by status, priority, category, date range

### 16.3 Testing Commands

```bash
# Run all tests
mix test

# Run with coverage
mix test --cover

# Run specific context
mix test test/tania/growth_test.exs

# Run in watch mode (with mix_test_watch)
mix test.watch
```

---

## 17. Migration Phases

### Phase 1: Project Bootstrap & Core Infrastructure (Week 1-2)

- [ ] Generate Phoenix project: `mix phx.new tania --binary-id --database postgres`
- [ ] Configure dual database support (PostgreSQL + SQLite)
- [ ] Set up Tailwind CSS with Tania brand colors (purple #7136A1, green #84BA4B)
- [ ] Generate authentication: `mix phx.gen.auth Accounts User users`
- [ ] Create shared layout (sidebar navigation matching current aside.vue, responsive design)
- [ ] Set up core_components (buttons, modals, tables, forms, pagination)
- [ ] Set up Leaflet map JS hook (for farm creation geolocation, replacing vue2-leaflet)
- [ ] Set up quantity slider JS hook (for crop move/dump, replacing vue-slider-component)
- [ ] Configure Gettext and migrate existing .po translation files from `resources/js/languages/`
- [ ] Set up Oban for background jobs
- [ ] Create AuditLog schema and helper module
- [ ] Write seeds.exs for default user
- [ ] Implement onboarding flow (ensure_has_farm on_mount hook + intro LiveViews)

### Phase 2: Farming Context (Week 2-3)

- [ ] Create Farm schema, migration, context
- [ ] Create Reservoir schema, migration, context
- [ ] Create Area schema, migration, context
- [ ] Create Note (polymorphic) schema, migration
- [ ] Build FarmLive.Index and FarmLive.Show LiveViews
- [ ] Build AreaLive.Index and AreaLive.Show LiveViews
- [ ] Build ReservoirLive.Index and ReservoirLive.Show LiveViews
- [ ] Implement notes CRUD for areas and reservoirs
- [ ] Implement area photo upload
- [ ] Write context tests and LiveView tests

### Phase 3: Inventory Context (Week 3-4)

- [ ] Create Material schema, migration, context
- [ ] Handle material types and embedded type_data
- [ ] Build MaterialLive.Index and MaterialLive.Show LiveViews
- [ ] Material filtering by type
- [ ] Plant type categorization
- [ ] Write tests

### Phase 4: Growth Context (Week 4-6)

- [ ] Create Crop schema and all sub-schemas (movement, harvest, dump, activity, photo)
- [ ] Create all migrations
- [ ] Implement crop batch creation with batch_id generation
- [ ] Implement crop movement logic (area transfers with quantity tracking)
- [ ] Implement harvest recording (full/partial with weight)
- [ ] Implement dump/discard recording
- [ ] Implement care activities (water, fertilize, prune, pesticide)
- [ ] Implement crop photo uploads
- [ ] Implement crop archival
- [ ] Build CropLive.Index (with active/archived tabs)
- [ ] Build CropLive.Show with all action modals (harvest, move, dump)
- [ ] Activity timeline view
- [ ] Days since seeding display
- [ ] Write extensive tests for crop lifecycle

### Phase 5: Tasks Context (Week 6-7)

- [ ] Create Task schema, migration, context
- [ ] Implement task lifecycle (create, complete, cancel)
- [ ] Asset linking (polymorphic association to area/crop/material/reservoir)
- [ ] Build TaskLive.Index with filtering (status, priority, category, date)
- [ ] Inline task creation from area/crop/reservoir detail pages
- [ ] Due date checking background job
- [ ] Write tests

### Phase 6: Dashboard & Polish (Week 7-8)

- [ ] Build DashboardLive with summary stats
- [ ] Active crops on production table
- [ ] Recent tasks widget
- [ ] PubSub real-time updates across all LiveViews
- [ ] Responsive design testing and fixes
- [ ] Account settings page (change password)
- [ ] Demo mode implementation
- [ ] Error pages (404, 500)
- [ ] Loading states and empty states

### Phase 7: JSON API (Week 8-9)

- [ ] Create API controllers for all resources
- [ ] API authentication (Bearer token)
- [ ] JSON serialization (Jason views or json_api library)
- [ ] API documentation (OpenAPI/Swagger)
- [ ] API tests

### Phase 8: Nerves Integration (Week 9-11)

- [ ] Set up Nerves project structure (poncho or conditional deps)
- [ ] Configure SQLite as database for Nerves target
- [ ] Configure VintageNet for WiFi/Ethernet
- [ ] Configure NTP time synchronization
- [ ] Set up data partition for database and uploads
- [ ] Replace Oban with Quantum for SQLite-compatible scheduling
- [ ] Build and test firmware for RPi 4
- [ ] Test firmware for RPi 3 and RPi Zero 2W
- [ ] Document firmware flashing process
- [ ] OTA update mechanism (optional)

### Phase 9: Deployment & Documentation (Week 11-12)

- [ ] Dockerfile and docker-compose.yml
- [ ] systemd service file
- [ ] Production configuration (runtime.exs)
- [ ] README with setup instructions
- [ ] Deployment guide (VPS)
- [ ] Nerves firmware guide (Raspberry Pi)
- [ ] Update Postman collection for new API
- [ ] CI/CD pipeline (GitHub Actions)

---

## 18. Dependencies & Libraries

### Core Dependencies

```elixir
# mix.exs
defp deps do
  [
    # Phoenix & Web
    {:phoenix, "~> 1.7"},
    {:phoenix_html, "~> 4.1"},
    {:phoenix_live_view, "~> 1.0"},
    {:phoenix_live_dashboard, "~> 0.8"},
    {:phoenix_live_reload, "~> 1.5", only: :dev},

    # Database
    {:ecto_sql, "~> 3.11"},
    {:postgrex, "~> 0.18"},          # PostgreSQL (VPS)
    {:ecto_sqlite3, "~> 0.16"},      # SQLite (Nerves + dev option)

    # Auth
    {:bcrypt_elixir, "~> 3.1"},

    # Background Jobs
    {:oban, "~> 2.17"},              # PostgreSQL targets
    {:quantum, "~> 3.5"},            # SQLite/Nerves targets

    # File Uploads
    {:waffle, "~> 1.1"},             # File upload abstraction
    {:waffle_ecto, "~> 0.0.12"},

    # Assets & Frontend
    {:tailwind, "~> 0.2", runtime: Mix.env() == :dev},
    {:esbuild, "~> 0.8", runtime: Mix.env() == :dev},
    {:heroicons, "~> 0.5"},          # Icon set

    # Utilities
    {:jason, "~> 1.4"},              # JSON encoding
    {:gettext, "~> 0.24"},           # i18n
    {:countries, "~> 1.6"},          # Country/city data
    {:telemetry_metrics, "~> 1.0"},
    {:telemetry_poller, "~> 1.1"},

    # Dev & Test
    {:floki, "~> 0.36", only: :test},
    {:mix_test_watch, "~> 1.2", only: :dev, runtime: false},
    {:credo, "~> 1.7", only: [:dev, :test], runtime: false},
    {:dialyxir, "~> 1.4", only: [:dev, :test], runtime: false},
  ]
end
```

### Nerves-Specific Dependencies

```elixir
# Only included when MIX_TARGET != "host"
{:nerves, "~> 1.10"},
{:nerves_system_rpi4, "~> 1.24", runtime: false, targets: :rpi4},
{:nerves_system_rpi3, "~> 1.25", runtime: false, targets: :rpi3},
{:nerves_system_rpi0_2, "~> 1.24", runtime: false, targets: :rpi0_2},
{:vintage_net, "~> 0.13"},
{:vintage_net_wifi, "~> 0.12"},
{:vintage_net_ethernet, "~> 0.11"},
{:nerves_time, "~> 0.4"},
{:ring_logger, "~> 0.10"},
{:nerves_pack, "~> 0.7"},
{:mdns_lite, "~> 0.8"},             # mDNS for tania.local
```

---

## 19. Configuration Management

### 19.1 Mapping from Go conf.json

| Go Config Key | Elixir Equivalent | Location |
|--------------|-------------------|----------|
| app_port | `PHX_PORT` env var | runtime.exs |
| tania_persistence_engine | Repo adapter selection | config.exs / runtime.exs |
| demo_mode | `config :tania, :demo_mode` | runtime.exs |
| upload_path_area | `config :tania, :upload_path` | config.exs |
| upload_path_crop | (same, organized by type) | config.exs |
| sqlite_path | `database` in Repo config | runtime.exs |
| mysql_* | Replaced by `DATABASE_URL` | runtime.exs |
| redirect_uri | Not needed (server-rendered) | - |
| client_id | Not needed (session-based) | - |

### 19.2 Environment Variables (Production)

```bash
DATABASE_URL=ecto://user:pass@host/tania_prod
SECRET_KEY_BASE=<64+ char random string>
PHX_HOST=tania.example.com
PHX_PORT=4000
TANIA_DEMO_MODE=false
TANIA_UPLOAD_PATH=/var/lib/tania/uploads
```

---

## 20. Open Questions & Decisions

### Must Decide Before Starting

1. **Project structure for Nerves**: Poncho project (separate mix projects sharing code) vs. single project with conditional compilation? **Recommendation**: Start with single project + conditional deps; split later if needed.

2. **Drop event sourcing completely or keep audit log?** **Recommendation**: Drop ES/CQRS, keep simple audit_log table.

3. **Email-based auth or username-based?** The current app uses username. Phoenix gen_auth defaults to email. **Recommendation**: Switch to email-based auth (modern standard).

4. **Multi-farm or single-farm per instance?** Current app supports multiple farms. Keep this? **Recommendation**: Keep multi-farm support.

5. **Multi-user roles?** Current app has single user type. Add roles (admin, worker, viewer)? **Recommendation**: Defer to v2. Keep single user type for now.

6. **Map integration**: The current Vue app uses Leaflet (vue2-leaflet) with OpenStreetMap tiles for farm geolocation. In LiveView, use a JS hook wrapping Leaflet directly. **Recommendation**: Implement as a Phoenix LiveView JS hook.

7. **Quantity sliders**: The Vue app uses vue-slider-component for crop move/dump quantities. **Recommendation**: Use a Leaflet-style JS hook wrapping a lightweight slider library, or use native HTML range input with LiveView event handling.

8. **PWA/Offline support**: The Vue app has SWPrecacheWebpackPlugin for service worker. **Recommendation**: Defer to later phase. LiveView requires connection anyway.

### Nice to Have (Future Phases)

- Weather API integration
- Sensor data from Nerves GPIO (soil moisture, temperature)
- Automated irrigation control via GPIO relays
- Mobile app (React Native or Flutter consuming the JSON API)
- Data export (CSV, PDF reports)
- Calendar view for tasks
- Notifications (email, push)
- Multi-tenant support (SaaS mode)
- Image recognition for plant health (ML integration)
- PWA/offline support with service worker

---

## Appendix A: Color Scheme

Preserve the Tania brand from the current SCSS:

```css
/* Tailwind custom colors in tailwind.config.js */
colors: {
  tania: {
    purple: '#7136A1',      /* Primary */
    green: '#84BA4B',       /* Secondary/Success */
    sidebar: '#513969',     /* Sidebar background */
    gray: '#F6F8F8',        /* Page background */
  }
}
```

## Appendix B: Existing Go Domain Enum Values

These must be preserved exactly for data compatibility:

**Farm Types**: organic, hydroponic, aquaponic, mushroom, livestock, fisheries, permaculture
**Area Types**: seeding, growing
**Area Locations**: indoor, outdoor
**Area Size Units**: m2 (sqm), Ha (hectare)
**Material Types**: seed, agrochemical, growing_medium, label_crop_support, seeding_container, post_harvest, tool, other
**Material Quantity Units**: seeds, packets, gram, kilogram, bags, bottles, cubic_metre, pieces, units
**Crop Status**: active, archived
**Crop Types**: seeding, growing
**Container Types**: tray, pot
**Task Status**: pending (created), completed, cancelled
**Task Priority**: high, medium, low
**Task Categories**: farm_maintenance, inventory, planting, harvesting, water_management, pest_control, safety, sanitation, general
**Task Domains**: area, crop, material, reservoir
**Water Source Types**: tap, well, rain_water, river, tank

## Appendix C: Vue.js to LiveView Component Mapping

| Vue Component | LiveView Equivalent | Notes |
|---------------|-------------------|-------|
| `app.vue` + `aside.vue` + `header.vue` | `layouts.ex` (root layout) | Sidebar + navbar in layout |
| `modal.vue` | `core_components.ex` (`.modal`) | Phoenix generator includes modal |
| `pagination.vue` | `core_components.ex` or Flop | Pagination component |
| `upload.vue` | LiveView `allow_upload/3` | Built-in upload support |
| `mapbox.vue` (Leaflet) | JS hook `leaflet_map.js` | Leaflet via phx-hook |
| `vue-slider-component` | JS hook or HTML range input | For crop move/dump quantities |
| `vuejs-datepicker` | Native HTML `<input type="date">` | Or JS hook if needed |
| `vue-toasted` | LiveView `put_flash/3` | Built-in flash messages |
| Vuex store modules | Phoenix Contexts | Business logic in contexts |
| Vue Router guards | LiveView `on_mount` hooks | Auth + onboarding checks |
| VeeValidate | Ecto changesets | Server-side validation |
| `vue-gettext` `<translate>` | `gettext()` / `ngettext()` | Same .po file format |
| Axios + auth interceptor | LiveView (no HTTP needed) | LiveView uses WebSocket |

## Appendix D: Existing Translation Keys

The existing .po files (in `resources/js/languages/`) contain approximately 200+ translatable strings covering:
- Navigation labels
- Form field labels and placeholders
- Button labels (Save, Cancel, Delete, etc.)
- Status labels (Active, Archived, Completed, etc.)
- Error messages
- Dashboard headings
- Table column headers
- Modal titles
- Onboarding wizard text
- Crop activity descriptions

These should be extracted and organized into Phoenix Gettext domains:
- `default.po` - General UI strings
- `errors.po` - Validation and error messages

## Appendix E: Key Files Reference (master branch)

| Purpose | Path (use `git show origin/master:<path>`) |
|---------|-----|
| Vue entry point | `resources/js/app.js` |
| Vue Router | `resources/js/router.js` |
| Vuex root store | `resources/js/stores/index.js` |
| All API endpoints | `resources/js/stores/api/farm.js` |
| Mutation types | `resources/js/stores/mutation-types.js` |
| Dashboard | `resources/js/pages/home.vue` |
| Crop detail | `resources/js/pages/farms/crop.vue` |
| Task page | `resources/js/pages/tasks/task.vue` |
| Onboarding farm | `resources/js/pages/intro/farm.vue` |
| HTTP service | `resources/js/services/http.js` |
| SCSS variables | `resources/sass/_variables.scss` |
| Webpack config | `webpack.config.js` |
| Go entry point | `backend/cmd/taniad/main.go` |
| Go config | `backend/conf.json` |
| SQLite schema | `backend/database/sqlite/ddl.sql` |
