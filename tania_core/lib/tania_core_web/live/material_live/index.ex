defmodule TaniaCoreWeb.MaterialLive.Index do
  use TaniaCoreWeb, :live_view

  alias TaniaCore.Farming
  alias TaniaCore.Inventory
  alias TaniaCore.Inventory.Material

  @impl true
  def mount(_params, _session, socket) do
    farm = Farming.get_user_farm(socket.assigns.current_scope.user.id)

    if farm do
      materials = Inventory.list_materials(farm.id)

      {:ok,
       socket
       |> assign(:page_title, "Materials")
       |> assign(:active_page, :materials)
       |> assign(:farm, farm)
       |> assign(:materials_count, length(materials))
       |> stream(:materials, materials)
       |> assign(:material, nil)
       |> assign(:show_form, false)
       |> assign(:quantity_units, [])}
    else
      {:ok, push_navigate(socket, to: ~p"/")}
    end
  end

  @impl true
  def handle_event("new_material", _params, socket) do
    {:noreply,
     socket
     |> assign(:material, %Material{})
     |> assign(:show_form, true)
     |> assign(:quantity_units, [])
     |> assign(:form, to_form(Inventory.change_material(%Material{})))}
  end

  def handle_event("close_form", _params, socket) do
    {:noreply, assign(socket, :show_form, false)}
  end

  def handle_event("validate", %{"material" => params}, socket) do
    material = socket.assigns.material || %Material{}
    changeset = Inventory.change_material(material, params)

    # Update quantity units when type changes
    quantity_units =
      case params["type"] do
        nil -> []
        "" -> []
        type -> Inventory.quantity_units_for_type(type)
      end

    {:noreply,
     socket
     |> assign(:form, to_form(Map.put(changeset, :action, :validate)))
     |> assign(:quantity_units, quantity_units)}
  end

  def handle_event("save", %{"material" => params}, socket) do
    save_material(socket, socket.assigns.material, params)
  end

  def handle_event("edit", %{"id" => id}, socket) do
    material = Inventory.get_material!(id)
    quantity_units = Inventory.quantity_units_for_type(material.type)

    {:noreply,
     socket
     |> assign(:material, material)
     |> assign(:show_form, true)
     |> assign(:quantity_units, quantity_units)
     |> assign(:form, to_form(Inventory.change_material(material)))}
  end

  def handle_event("delete", %{"id" => id}, socket) do
    material = Inventory.get_material!(id)
    user_id = socket.assigns.current_scope.user.id

    case Inventory.delete_material(material, user_id) do
      {:ok, _} ->
        {:noreply,
         socket
         |> stream_delete(:materials, material)
         |> update(:materials_count, &(&1 - 1))
         |> put_flash(:info, "Material deleted")}

      {:error, _} ->
        {:noreply, put_flash(socket, :error, "Could not delete material")}
    end
  end

  defp save_material(socket, %Material{id: nil}, params) do
    params = Map.put(params, "farm_id", socket.assigns.farm.id)
    user_id = socket.assigns.current_scope.user.id

    case Inventory.create_material(params, user_id) do
      {:ok, material} ->
        {:noreply,
         socket
         |> stream_insert(:materials, material)
         |> update(:materials_count, &(&1 + 1))
         |> assign(:show_form, false)
         |> put_flash(:info, "Material created")}

      {:error, changeset} ->
        {:noreply, assign(socket, :form, to_form(changeset))}
    end
  end

  defp save_material(socket, material, params) do
    user_id = socket.assigns.current_scope.user.id

    case Inventory.update_material(material, params, user_id) do
      {:ok, material} ->
        {:noreply,
         socket
         |> stream_insert(:materials, material)
         |> assign(:show_form, false)
         |> put_flash(:info, "Material updated")}

      {:error, changeset} ->
        {:noreply, assign(socket, :form, to_form(changeset))}
    end
  end

  @type_labels %{
    "seed" => "Seed",
    "agrochemical" => "Agrochemical",
    "growing_medium" => "Growing Medium",
    "label_crop_support" => "Label & Crop Support",
    "seeding_container" => "Seeding Container",
    "post_harvest_supply" => "Post Harvest Supply",
    "plant" => "Plant",
    "other" => "Other"
  }

  defp type_label(type), do: Map.get(@type_labels, type, type)

  defp type_options do
    Enum.map(@type_labels, fn {key, label} -> {label, key} end)
    |> Enum.sort_by(&elem(&1, 0))
  end

  defp unit_options(units) do
    Enum.map(units, fn unit ->
      label = unit |> String.replace("_", " ") |> String.capitalize()
      {label, unit}
    end)
  end

  @impl true
  def render(assigns) do
    ~H"""
    <Layouts.app flash={@flash} current_scope={@current_scope} active_page={:materials}>
      <.header>
        {gettext("Materials")}
        <:actions>
          <.button phx-click="new_material">
            <.icon name="hero-plus" class="size-4" /> {gettext("New Material")}
          </.button>
        </:actions>
      </.header>

      <div class="mt-6">
        <.table
          :if={@materials_count > 0}
          id="materials"
          rows={@streams.materials}
        >
          <:col :let={{_id, m}} label={gettext("Name")}>{m.name}</:col>
          <:col :let={{_id, m}} label={gettext("Type")}>
            <.badge variant="primary">{type_label(m.type)}</.badge>
          </:col>
          <:col :let={{_id, m}} label={gettext("Quantity")}>
            {m.quantity_value} {m.quantity_unit}
          </:col>
          <:col :let={{_id, m}} label={gettext("Price")}>
            {if m.price_per_unit, do: "#{m.price_per_unit} #{m.price_currency}", else: "-"}
          </:col>
          <:action :let={{_id, m}}>
            <.button variant="ghost" phx-click={JS.push("edit", value: %{id: m.id})}>
              <.icon name="hero-pencil-square" class="size-4" />
            </.button>
            <.button
              variant="ghost"
              phx-click={JS.push("delete", value: %{id: m.id})}
              data-confirm="Are you sure?"
            >
              <.icon name="hero-trash" class="size-4" />
            </.button>
          </:action>
        </.table>

        <.empty_state
          :if={@materials_count == 0}
          icon="hero-archive-box"
          message={gettext("No materials yet. Add seeds, chemicals, or other supplies.")}
        >
          <:actions>
            <.button phx-click="new_material">
              <.icon name="hero-plus" class="size-4" /> {gettext("New Material")}
            </.button>
          </:actions>
        </.empty_state>
      </div>

      <.modal :if={@show_form} id="material-form-modal" show on_cancel={JS.push("close_form")}>
        <h2 class="text-lg font-semibold mb-4">
          {if @material && @material.id, do: gettext("Edit Material"), else: gettext("New Material")}
        </h2>
        <.form for={@form} phx-change="validate" phx-submit="save">
          <.input field={@form[:name]} label={gettext("Name")} required />
          <.input
            field={@form[:type]}
            type="select"
            label={gettext("Type")}
            options={type_options()}
            prompt="Select type"
            required
          />

          <div :if={@form[:type].value in ["seed", "plant"]} class="mb-3">
            <.input
              field={@form[:type_data]}
              type="select"
              label={gettext("Plant Type")}
              name="material[type_data][plant_type]"
              options={Enum.map(Material.plant_types(), &{String.capitalize(&1), &1})}
              prompt="Select plant type"
              value={
                get_in(@form.source.changes, [:type_data, "plant_type"]) ||
                  get_in(@form.data.type_data || %{}, ["plant_type"])
              }
            />
          </div>

          <div :if={@form[:type].value == "agrochemical"} class="mb-3">
            <.input
              field={@form[:type_data]}
              type="select"
              label={gettext("Chemical Type")}
              name="material[type_data][chemical_type]"
              options={Enum.map(Material.chemical_types(), &{String.capitalize(&1), &1})}
              prompt="Select chemical type"
              value={
                get_in(@form.source.changes, [:type_data, "chemical_type"]) ||
                  get_in(@form.data.type_data || %{}, ["chemical_type"])
              }
            />
          </div>

          <div :if={@form[:type].value == "seeding_container"} class="mb-3">
            <.input
              field={@form[:type_data]}
              type="select"
              label={gettext("Container Type")}
              name="material[type_data][container_type]"
              options={Enum.map(Material.container_types(), &{String.capitalize(&1), &1})}
              prompt="Select container type"
              value={
                get_in(@form.source.changes, [:type_data, "container_type"]) ||
                  get_in(@form.data.type_data || %{}, ["container_type"])
              }
            />
          </div>

          <.input field={@form[:quantity_value]} type="number" label={gettext("Quantity")} step="any" />
          <.input
            :if={@quantity_units != []}
            field={@form[:quantity_unit]}
            type="select"
            label={gettext("Unit")}
            options={unit_options(@quantity_units)}
            prompt="Select unit"
          />
          <.input
            field={@form[:price_per_unit]}
            type="number"
            label={gettext("Price per Unit")}
            step="any"
          />
          <.input field={@form[:price_currency]} label={gettext("Currency")} />
          <.input field={@form[:expiration_date]} type="date" label={gettext("Expiration Date")} />
          <.input field={@form[:produced_by]} label={gettext("Produced By")} />
          <.input field={@form[:notes]} type="textarea" label={gettext("Notes")} />

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
