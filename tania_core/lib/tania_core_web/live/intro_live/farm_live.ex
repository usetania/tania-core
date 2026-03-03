defmodule TaniaCoreWeb.IntroLive.FarmLive do
  use TaniaCoreWeb, :live_view

  alias TaniaCore.Farming
  alias TaniaCore.Farming.Farm

  @impl true
  def mount(_params, _session, socket) do
    changeset = Farming.change_farm(%Farm{})

    {:ok,
     socket
     |> assign(:page_title, "Create Your Farm")
     |> assign(:form, to_form(changeset))
     |> assign(:step, 1)}
  end

  @impl true
  def handle_event("validate", %{"farm" => farm_params}, socket) do
    changeset =
      %Farm{}
      |> Farming.change_farm(farm_params)
      |> Map.put(:action, :validate)

    {:noreply, assign(socket, :form, to_form(changeset))}
  end

  def handle_event("save", %{"farm" => farm_params}, socket) do
    user_id = socket.assigns.current_scope.user.id

    case Farming.create_farm(farm_params, user_id) do
      {:ok, _farm} ->
        {:noreply, push_navigate(socket, to: "/intro/reservoir")}

      {:error, changeset} ->
        {:noreply, assign(socket, :form, to_form(changeset))}
    end
  end

  def handle_event("map_clicked", %{"lat" => lat, "lng" => lng}, socket) do
    form = socket.assigns.form

    params =
      form.source.changes
      |> Map.put(:latitude, lat)
      |> Map.put(:longitude, lng)
      |> Enum.into(%{}, fn {k, v} -> {to_string(k), v} end)

    changeset = Farming.change_farm(%Farm{}, params)
    {:noreply, assign(socket, :form, to_form(changeset))}
  end

  @farm_types [
    {"Organic", "organic"},
    {"Conventional", "conventional"},
    {"Hydroponic", "hydroponic"},
    {"Aquaponic", "aquaponic"},
    {"Nursery", "nursery"},
    {"Urban", "urban"},
    {"Permaculture", "permaculture"}
  ]

  @impl true
  def render(assigns) do
    assigns = assign(assigns, :farm_types, @farm_types)

    ~H"""
    <div class="min-h-screen bg-background flex items-center justify-center p-4">
      <div class="w-full max-w-lg">
        <div class="text-center mb-8">
          <div class="flex justify-center mb-4">
            <div class="flex size-12 items-center justify-center rounded-xl bg-primary">
              <.icon name="hero-sun" class="size-7 text-white" />
            </div>
          </div>
          <h1 class="text-2xl font-bold text-text">{gettext("Welcome to Tania")}</h1>
          <p class="text-text-light mt-2">
            {gettext("Let's set up your farm. Step 1 of 3.")}
          </p>
        </div>

        <%!-- Progress bar --%>
        <div class="flex gap-2 mb-8">
          <div class="h-1 flex-1 rounded-full bg-primary" />
          <div class="h-1 flex-1 rounded-full bg-border" />
          <div class="h-1 flex-1 rounded-full bg-border" />
        </div>

        <div class="rounded-xl border border-border bg-white p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-text mb-4">{gettext("Create Your Farm")}</h2>
          <.form for={@form} phx-change="validate" phx-submit="save">
            <.input field={@form[:name]} label={gettext("Farm Name")} required />
            <.input
              field={@form[:type]}
              type="select"
              label={gettext("Farm Type")}
              options={@farm_types}
              prompt={gettext("Select farm type")}
              required
            />
            <.input field={@form[:country]} label={gettext("Country")} />
            <.input field={@form[:city]} label={gettext("City")} />

            <div class="mb-3">
              <label class="block text-sm font-medium text-text mb-1">
                {gettext("Location (click the map)")}
              </label>
              <div
                id="farm-map"
                phx-hook="LeafletMap"
                phx-update="ignore"
                data-lat={@form[:latitude].value || "0"}
                data-lng={@form[:longitude].value || "0"}
                data-zoom="3"
                class="h-64 rounded-lg border border-border"
              />
            </div>

            <.input field={@form[:latitude]} type="number" label={gettext("Latitude")} step="any" />
            <.input field={@form[:longitude]} type="number" label={gettext("Longitude")} step="any" />

            <div class="mt-6">
              <.button type="submit" class="w-full">
                {gettext("Continue")} <.icon name="hero-arrow-right" class="size-4" />
              </.button>
            </div>
          </.form>
        </div>
      </div>
    </div>
    """
  end
end
