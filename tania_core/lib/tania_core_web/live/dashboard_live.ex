defmodule TaniaCoreWeb.DashboardLive do
  use TaniaCoreWeb, :live_view

  alias TaniaCore.Farming
  alias TaniaCore.Growth
  alias TaniaCore.Tasks

  @impl true
  def mount(_params, _session, socket) do
    farm = Farming.get_user_farm(socket.assigns.current_scope.user.id)

    if farm do
      areas = Farming.list_areas(farm.id)
      active_crops = Growth.list_crops(farm.id, "active")

      seeding_count = Enum.count(active_crops, &(&1.type == "seeding"))
      growing_count = Enum.count(active_crops, &(&1.type == "growing"))

      recent_tasks = Tasks.list_tasks(farm.id, status: "created") |> Enum.take(5)

      {:ok,
       socket
       |> assign(:page_title, "Dashboard")
       |> assign(:active_page, :dashboard)
       |> assign(:farm, farm)
       |> assign(:area_count, length(areas))
       |> assign(:active_crop_count, length(active_crops))
       |> assign(:seeding_count, seeding_count)
       |> assign(:growing_count, growing_count)
       |> assign(:recent_crops, Enum.take(active_crops, 10))
       |> assign(:recent_tasks, recent_tasks)}
    else
      {:ok, push_navigate(socket, to: ~p"/")}
    end
  end

  defp format_date(nil), do: "-"
  defp format_date(dt), do: Calendar.strftime(dt, "%b %d, %Y")

  defp priority_variant("urgent"), do: "danger"
  defp priority_variant(_), do: "default"

  @impl true
  def render(assigns) do
    ~H"""
    <Layouts.app flash={@flash} current_scope={@current_scope} active_page={:dashboard}>
      <.header>
        {gettext("Dashboard")}
        <:subtitle>{gettext("Welcome to")} {@farm.name}</:subtitle>
      </.header>

      <%!-- Stats cards --%>
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-6">
        <.stat_card label={gettext("Areas")} value={@area_count} icon="hero-map" />
        <.stat_card
          label={gettext("Active Crops")}
          value={@active_crop_count}
          icon="hero-squares-2x2"
        />
        <.stat_card label={gettext("Seeding")} value={@seeding_count} icon="hero-sun" />
        <.stat_card label={gettext("Growing")} value={@growing_count} icon="hero-arrow-trending-up" />
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-8">
        <%!-- Active crops table --%>
        <.card>
          <:header>{gettext("Active Crops")}</:header>
          <.table :if={@recent_crops != []} id="dashboard-crops" rows={@recent_crops}>
            <:col :let={crop} label={gettext("Batch ID")}>
              <.link
                navigate={"/crops/#{crop.id}"}
                class="text-primary hover:underline font-mono text-xs"
              >
                {crop.batch_id}
              </.link>
            </:col>
            <:col :let={crop} label={gettext("Type")}>
              <.badge variant={if(crop.type == "seeding", do: "info", else: "success")}>
                {String.capitalize(crop.type)}
              </.badge>
            </:col>
            <:col :let={crop} label={gettext("Qty")}>{crop.quantity}</:col>
          </.table>
          <p :if={@recent_crops == []} class="text-sm text-text-light text-center py-4">
            {gettext("No active crops.")}
          </p>
        </.card>

        <%!-- Recent tasks --%>
        <.card>
          <:header>{gettext("Recent Tasks")}</:header>
          <div :if={@recent_tasks != []} class="divide-y divide-border">
            <div :for={task <- @recent_tasks} class="py-3 first:pt-0 last:pb-0">
              <div class="flex items-center justify-between">
                <span class="text-sm text-text font-medium">{task.title}</span>
                <.badge variant={priority_variant(task.priority)}>
                  {String.capitalize(task.priority)}
                </.badge>
              </div>
              <p :if={task.due_date} class="text-xs text-text-light mt-1">
                {gettext("Due")}: {format_date(task.due_date)}
              </p>
            </div>
          </div>
          <p :if={@recent_tasks == []} class="text-sm text-text-light text-center py-4">
            {gettext("No pending tasks.")}
          </p>
        </.card>
      </div>
    </Layouts.app>
    """
  end
end
