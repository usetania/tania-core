defmodule TaniaCoreWeb.AreaLive.Index do
  use TaniaCoreWeb, :live_view

  alias TaniaCore.Farming
  alias TaniaCore.Farming.Area

  @impl true
  def mount(_params, _session, socket) do
    farm = Farming.get_user_farm(socket.assigns.current_scope.user.id)

    if farm do
      areas = Farming.list_areas(farm.id)
      reservoirs = Farming.list_reservoirs(farm.id)

      {:ok,
       socket
       |> assign(:page_title, "Areas")
       |> assign(:active_page, :areas)
       |> assign(:farm, farm)
       |> assign(:reservoirs, reservoirs)
       |> stream(:areas, areas)
       |> assign(:area, nil)
       |> assign(:show_form, false)}
    else
      {:ok, push_navigate(socket, to: ~p"/")}
    end
  end

  @impl true
  def handle_params(params, _url, socket) do
    {:noreply, apply_action(socket, socket.assigns.live_action, params)}
  end

  defp apply_action(socket, :index, _params) do
    assign(socket, :area, nil)
  end

  defp apply_action(socket, :new, _params) do
    socket
    |> assign(:area, %Area{})
    |> assign(:show_form, true)
    |> assign(:form, to_form(Farming.change_area(%Area{})))
  end

  defp apply_action(socket, :edit, %{"id" => id}) do
    area = Farming.get_area!(id)

    socket
    |> assign(:area, area)
    |> assign(:show_form, true)
    |> assign(:form, to_form(Farming.change_area(area)))
  end

  @impl true
  def handle_event("new_area", _params, socket) do
    {:noreply,
     socket
     |> assign(:area, %Area{})
     |> assign(:show_form, true)
     |> assign(:form, to_form(Farming.change_area(%Area{})))}
  end

  def handle_event("close_form", _params, socket) do
    {:noreply, assign(socket, :show_form, false)}
  end

  def handle_event("validate", %{"area" => area_params}, socket) do
    area = socket.assigns.area || %Area{}
    changeset = Farming.change_area(area, area_params)
    {:noreply, assign(socket, :form, to_form(Map.put(changeset, :action, :validate)))}
  end

  def handle_event("save", %{"area" => area_params}, socket) do
    save_area(socket, socket.assigns.area, area_params)
  end

  def handle_event("delete", %{"id" => id}, socket) do
    area = Farming.get_area!(id)
    user_id = socket.assigns.current_scope.user.id

    case Farming.delete_area(area, user_id) do
      {:ok, _} ->
        {:noreply,
         socket
         |> stream_delete(:areas, area)
         |> put_flash(:info, "Area deleted")}

      {:error, _} ->
        {:noreply, put_flash(socket, :error, "Could not delete area")}
    end
  end

  defp save_area(socket, %Area{id: nil}, area_params) do
    params = Map.put(area_params, "farm_id", socket.assigns.farm.id)
    user_id = socket.assigns.current_scope.user.id

    case Farming.create_area(params, user_id) do
      {:ok, area} ->
        area = Farming.get_area!(area.id)

        {:noreply,
         socket
         |> stream_insert(:areas, area)
         |> assign(:show_form, false)
         |> put_flash(:info, "Area created")}

      {:error, changeset} ->
        {:noreply, assign(socket, :form, to_form(changeset))}
    end
  end

  defp save_area(socket, area, area_params) do
    user_id = socket.assigns.current_scope.user.id

    case Farming.update_area(area, area_params, user_id) do
      {:ok, area} ->
        area = Farming.get_area!(area.id)

        {:noreply,
         socket
         |> stream_insert(:areas, area)
         |> assign(:show_form, false)
         |> put_flash(:info, "Area updated")}

      {:error, changeset} ->
        {:noreply, assign(socket, :form, to_form(changeset))}
    end
  end

  @impl true
  def render(assigns) do
    ~H"""
    <Layouts.app flash={@flash} current_scope={@current_scope} active_page={:areas}>
      <.header>
        {gettext("Areas")}
        <:actions>
          <.button phx-click="new_area">
            <.icon name="hero-plus" class="size-4" /> {gettext("New Area")}
          </.button>
        </:actions>
      </.header>

      <div
        :if={@streams.areas |> Enum.any?()}
        class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 mt-6"
      >
        <div
          :for={{dom_id, area} <- @streams.areas}
          id={dom_id}
          class="rounded-xl border border-border bg-white shadow-sm overflow-hidden hover:shadow-md transition-shadow"
        >
          <div class="aspect-video bg-gray-100 flex items-center justify-center">
            <img :if={area.photo_url} src={area.photo_url} class="h-full w-full object-cover" />
            <.icon :if={!area.photo_url} name="hero-photo" class="size-12 text-text-light/30" />
          </div>
          <div class="p-4">
            <h3 class="font-semibold text-text">{area.name}</h3>
            <div class="mt-2 flex flex-wrap gap-2">
              <.badge variant={if(area.type == "seeding", do: "info", else: "success")}>
                {String.capitalize(area.type)}
              </.badge>
              <.badge>{String.capitalize(area.location)}</.badge>
              <.badge :if={area.size_value} variant="default">
                {area.size_value} {area.size_unit}
              </.badge>
            </div>
            <p :if={area.reservoir} class="mt-2 text-xs text-text-light">
              {gettext("Reservoir")}: {area.reservoir.name}
            </p>
            <div class="mt-3 flex gap-2">
              <.button
                variant="ghost"
                phx-click={JS.push("delete", value: %{id: area.id})}
                data-confirm="Are you sure?"
              >
                <.icon name="hero-trash" class="size-4" />
              </.button>
            </div>
          </div>
        </div>
      </div>

      <.empty_state
        :if={!(@streams.areas |> Enum.any?())}
        icon="hero-map"
        message={gettext("No areas yet. Create your first area to get started.")}
      >
        <:actions>
          <.button phx-click="new_area">
            <.icon name="hero-plus" class="size-4" /> {gettext("New Area")}
          </.button>
        </:actions>
      </.empty_state>

      <.modal :if={@show_form} id="area-form-modal" show on_cancel={JS.push("close_form")}>
        <h2 class="text-lg font-semibold mb-4">
          {if @area && @area.id, do: gettext("Edit Area"), else: gettext("New Area")}
        </h2>
        <.form for={@form} phx-change="validate" phx-submit="save">
          <.input field={@form[:name]} label={gettext("Name")} required />
          <.input
            field={@form[:type]}
            type="select"
            label={gettext("Type")}
            options={[{"Seeding", "seeding"}, {"Growing", "growing"}]}
            required
          />
          <.input
            field={@form[:location]}
            type="select"
            label={gettext("Location")}
            options={[{"Field (Outdoor)", "outdoor"}, {"Greenhouse (Indoor)", "indoor"}]}
            required
          />
          <.input field={@form[:size_value]} type="number" label={gettext("Size")} step="any" />
          <.input
            field={@form[:size_unit]}
            type="select"
            label={gettext("Unit")}
            options={[{"Square Meter", "m2"}, {"Hectare", "Ha"}]}
            prompt="Select unit"
          />
          <.input
            field={@form[:reservoir_id]}
            type="select"
            label={gettext("Reservoir")}
            options={Enum.map(@reservoirs, &{&1.name, &1.id})}
            prompt="Select reservoir (optional)"
          />
          <div class="mt-4 flex justify-end gap-2">
            <.button type="button" variant="secondary" phx-click="close_form">
              {gettext("Cancel")}
            </.button>
            <.button type="submit">{gettext("Save")}</.button>
          </div>
        </.form>
      </.modal>
    </Layouts.app>
    """
  end
end
