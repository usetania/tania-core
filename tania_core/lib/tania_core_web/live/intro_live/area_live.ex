defmodule TaniaCoreWeb.IntroLive.AreaLive do
  use TaniaCoreWeb, :live_view

  alias TaniaCore.Farming
  alias TaniaCore.Farming.Area

  @impl true
  def mount(_params, _session, socket) do
    farm = Farming.get_user_farm(socket.assigns.current_scope.user.id)

    if farm do
      reservoirs = Farming.list_reservoirs(farm.id)
      changeset = Farming.change_area(%Area{})

      {:ok,
       socket
       |> assign(:page_title, "Create Your First Area")
       |> assign(:farm, farm)
       |> assign(:reservoirs, reservoirs)
       |> assign(:form, to_form(changeset))}
    else
      {:ok, push_navigate(socket, to: "/intro/farm")}
    end
  end

  @impl true
  def handle_event("validate", %{"area" => params}, socket) do
    changeset =
      %Area{}
      |> Farming.change_area(params)
      |> Map.put(:action, :validate)

    {:noreply, assign(socket, :form, to_form(changeset))}
  end

  def handle_event("save", %{"area" => params}, socket) do
    params = Map.put(params, "farm_id", socket.assigns.farm.id)
    user_id = socket.assigns.current_scope.user.id

    case Farming.create_area(params, user_id) do
      {:ok, _area} ->
        {:noreply,
         socket
         |> put_flash(:info, gettext("Farm setup complete! Welcome to Tania."))
         |> push_navigate(to: "/dashboard")}

      {:error, changeset} ->
        {:noreply, assign(socket, :form, to_form(changeset))}
    end
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-background flex items-center justify-center p-4">
      <div class="w-full max-w-lg">
        <div class="text-center mb-8">
          <div class="flex justify-center mb-4">
            <div class="flex size-12 items-center justify-center rounded-xl bg-primary">
              <.icon name="hero-sun" class="size-7 text-white" />
            </div>
          </div>
          <h1 class="text-2xl font-bold text-text">{gettext("Create Your First Area")}</h1>
          <p class="text-text-light mt-2">
            {gettext("Step 3 of 3. Add a growing or seeding area.")}
          </p>
        </div>

        <%!-- Progress bar --%>
        <div class="flex gap-2 mb-8">
          <div class="h-1 flex-1 rounded-full bg-primary" />
          <div class="h-1 flex-1 rounded-full bg-primary" />
          <div class="h-1 flex-1 rounded-full bg-primary" />
        </div>

        <div class="rounded-xl border border-border bg-white p-6 shadow-sm">
          <.form for={@form} phx-change="validate" phx-submit="save">
            <.input field={@form[:name]} label={gettext("Area Name")} required />
            <.input
              field={@form[:type]}
              type="select"
              label={gettext("Type")}
              options={[
                {gettext("Seeding"), "seeding"},
                {gettext("Growing"), "growing"}
              ]}
              required
            />
            <.input
              field={@form[:location]}
              type="select"
              label={gettext("Location")}
              options={[
                {gettext("Field (Outdoor)"), "outdoor"},
                {gettext("Greenhouse (Indoor)"), "indoor"}
              ]}
              required
            />
            <.input field={@form[:size_value]} type="number" label={gettext("Size")} step="any" />
            <.input
              field={@form[:size_unit]}
              type="select"
              label={gettext("Unit")}
              options={[{"m\u00B2", "m2"}, {"Ha", "Ha"}]}
              prompt={gettext("Select unit")}
            />
            <.input
              :if={@reservoirs != []}
              field={@form[:reservoir_id]}
              type="select"
              label={gettext("Reservoir")}
              options={Enum.map(@reservoirs, &{&1.name, &1.id})}
              prompt={gettext("Select reservoir (optional)")}
            />

            <div class="mt-6">
              <.button type="submit" class="w-full">
                {gettext("Finish Setup")} <.icon name="hero-check" class="size-4" />
              </.button>
            </div>
          </.form>
        </div>
      </div>
    </div>
    """
  end
end
