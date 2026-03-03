defmodule TaniaCoreWeb.ReservoirLive.Index do
  use TaniaCoreWeb, :live_view

  alias TaniaCore.Farming
  alias TaniaCore.Farming.Reservoir

  @impl true
  def mount(_params, _session, socket) do
    farm = Farming.get_user_farm(socket.assigns.current_scope.user.id)

    if farm do
      reservoirs = Farming.list_reservoirs(farm.id)

      {:ok,
       socket
       |> assign(:page_title, "Reservoirs")
       |> assign(:active_page, :reservoirs)
       |> assign(:farm, farm)
       |> stream(:reservoirs, reservoirs)
       |> assign(:reservoir, nil)
       |> assign(:show_form, false)}
    else
      {:ok, push_navigate(socket, to: ~p"/")}
    end
  end

  @impl true
  def handle_event("new_reservoir", _params, socket) do
    {:noreply,
     socket
     |> assign(:reservoir, %Reservoir{})
     |> assign(:show_form, true)
     |> assign(:form, to_form(Farming.change_reservoir(%Reservoir{})))}
  end

  def handle_event("edit", %{"id" => id}, socket) do
    reservoir = Farming.get_reservoir!(id)

    {:noreply,
     socket
     |> assign(:reservoir, reservoir)
     |> assign(:show_form, true)
     |> assign(:form, to_form(Farming.change_reservoir(reservoir)))}
  end

  def handle_event("close_form", _params, socket) do
    {:noreply, assign(socket, :show_form, false)}
  end

  def handle_event("validate", %{"reservoir" => params}, socket) do
    reservoir = socket.assigns.reservoir || %Reservoir{}
    changeset = Farming.change_reservoir(reservoir, params)
    {:noreply, assign(socket, :form, to_form(Map.put(changeset, :action, :validate)))}
  end

  def handle_event("save", %{"reservoir" => params}, socket) do
    save_reservoir(socket, socket.assigns.reservoir, params)
  end

  def handle_event("delete", %{"id" => id}, socket) do
    reservoir = Farming.get_reservoir!(id)
    user_id = socket.assigns.current_scope.user.id

    case Farming.delete_reservoir(reservoir, user_id) do
      {:ok, _} ->
        {:noreply,
         socket
         |> stream_delete(:reservoirs, reservoir)
         |> put_flash(:info, "Reservoir deleted")}

      {:error, _} ->
        {:noreply, put_flash(socket, :error, "Could not delete reservoir")}
    end
  end

  defp save_reservoir(socket, %Reservoir{id: nil}, params) do
    params = Map.put(params, "farm_id", socket.assigns.farm.id)
    user_id = socket.assigns.current_scope.user.id

    case Farming.create_reservoir(params, user_id) do
      {:ok, reservoir} ->
        {:noreply,
         socket
         |> stream_insert(:reservoirs, reservoir)
         |> assign(:show_form, false)
         |> put_flash(:info, "Reservoir created")}

      {:error, changeset} ->
        {:noreply, assign(socket, :form, to_form(changeset))}
    end
  end

  defp save_reservoir(socket, reservoir, params) do
    user_id = socket.assigns.current_scope.user.id

    case Farming.update_reservoir(reservoir, params, user_id) do
      {:ok, reservoir} ->
        {:noreply,
         socket
         |> stream_insert(:reservoirs, reservoir)
         |> assign(:show_form, false)
         |> put_flash(:info, "Reservoir updated")}

      {:error, changeset} ->
        {:noreply, assign(socket, :form, to_form(changeset))}
    end
  end

  @impl true
  def render(assigns) do
    ~H"""
    <Layouts.app flash={@flash} current_scope={@current_scope} active_page={:reservoirs}>
      <.header>
        {gettext("Reservoirs")}
        <:actions>
          <.button phx-click="new_reservoir">
            <.icon name="hero-plus" class="size-4" /> {gettext("New Reservoir")}
          </.button>
        </:actions>
      </.header>

      <div class="mt-6">
        <.table
          :if={@streams.reservoirs |> Enum.any?()}
          id="reservoirs"
          rows={@streams.reservoirs}
        >
          <:col :let={{_id, reservoir}} label={gettext("Name")}>{reservoir.name}</:col>
          <:col :let={{_id, reservoir}} label={gettext("Water Source")}>
            {String.capitalize(reservoir.water_source_type)}
          </:col>
          <:col :let={{_id, reservoir}} label={gettext("Capacity")}>
            {if reservoir.capacity, do: "#{reservoir.capacity} L", else: "-"}
          </:col>
          <:action :let={{_id, reservoir}}>
            <.button variant="ghost" phx-click={JS.push("edit", value: %{id: reservoir.id})}>
              <.icon name="hero-pencil-square" class="size-4" />
            </.button>
            <.button
              variant="ghost"
              phx-click={JS.push("delete", value: %{id: reservoir.id})}
              data-confirm="Are you sure?"
            >
              <.icon name="hero-trash" class="size-4" />
            </.button>
          </:action>
        </.table>

        <.empty_state
          :if={!(@streams.reservoirs |> Enum.any?())}
          icon="hero-beaker"
          message={gettext("No reservoirs yet. Create your first water source.")}
        >
          <:actions>
            <.button phx-click="new_reservoir">
              <.icon name="hero-plus" class="size-4" /> {gettext("New Reservoir")}
            </.button>
          </:actions>
        </.empty_state>
      </div>

      <.modal :if={@show_form} id="reservoir-form-modal" show on_cancel={JS.push("close_form")}>
        <h2 class="text-lg font-semibold mb-4">
          {if @reservoir && @reservoir.id,
            do: gettext("Edit Reservoir"),
            else: gettext("New Reservoir")}
        </h2>
        <.form for={@form} phx-change="validate" phx-submit="save">
          <.input field={@form[:name]} label={gettext("Name")} required />
          <.input
            field={@form[:water_source_type]}
            type="select"
            label={gettext("Water Source Type")}
            options={[{"Tap / Well", "tap"}, {"Water Tank / Cistern", "bucket"}]}
            required
          />
          <.input
            field={@form[:capacity]}
            type="number"
            label={gettext("Capacity (Litres)")}
            step="any"
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
