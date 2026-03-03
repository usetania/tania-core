defmodule TaniaCoreWeb.TaskLive.Index do
  use TaniaCoreWeb, :live_view

  alias TaniaCore.Farming
  alias TaniaCore.Tasks
  alias TaniaCore.Tasks.Task

  @impl true
  def mount(_params, _session, socket) do
    farm = Farming.get_user_farm(socket.assigns.current_scope.user.id)

    if farm do
      tasks = Tasks.list_tasks(farm.id)

      {:ok,
       socket
       |> assign(:page_title, "Tasks")
       |> assign(:active_page, :tasks)
       |> assign(:farm, farm)
       |> assign(:filter, "all")
       |> stream(:tasks, tasks)
       |> assign(:task, nil)
       |> assign(:show_form, false)}
    else
      {:ok, push_navigate(socket, to: ~p"/")}
    end
  end

  @impl true
  def handle_event("filter", %{"filter" => filter}, socket) do
    farm_id = socket.assigns.farm.id

    tasks =
      case filter do
        "all" -> Tasks.list_tasks(farm_id)
        "incomplete" -> Tasks.list_tasks(farm_id, status: "created")
        "completed" -> Tasks.list_tasks(farm_id, status: "completed")
        "overdue" -> Tasks.list_tasks(farm_id, is_due: true)
        "today" -> Tasks.list_tasks(farm_id, due_today: true)
        "this_week" -> Tasks.list_tasks(farm_id, due_this_week: true)
        "this_month" -> Tasks.list_tasks(farm_id, due_this_month: true)
        _ -> Tasks.list_tasks(farm_id)
      end

    {:noreply,
     socket
     |> assign(:filter, filter)
     |> stream(:tasks, tasks, reset: true)}
  end

  def handle_event("new_task", _params, socket) do
    {:noreply,
     socket
     |> assign(:task, %Task{})
     |> assign(:show_form, true)
     |> assign(:form, to_form(Tasks.change_task(%Task{})))}
  end

  def handle_event("close_form", _params, socket) do
    {:noreply, assign(socket, :show_form, false)}
  end

  def handle_event("validate", %{"task" => task_params}, socket) do
    task = socket.assigns.task || %Task{}
    changeset = Tasks.change_task(task, task_params)
    {:noreply, assign(socket, :form, to_form(Map.put(changeset, :action, :validate)))}
  end

  def handle_event("save", %{"task" => task_params}, socket) do
    params = Map.put(task_params, "farm_id", socket.assigns.farm.id)
    user_id = socket.assigns.current_scope.user.id

    case Tasks.create_task(params, user_id) do
      {:ok, task} ->
        {:noreply,
         socket
         |> stream_insert(:tasks, task)
         |> assign(:show_form, false)
         |> put_flash(:info, "Task created")}

      {:error, changeset} ->
        {:noreply, assign(socket, :form, to_form(changeset))}
    end
  end

  def handle_event("complete", %{"id" => id}, socket) do
    task = Tasks.get_task!(id)
    user_id = socket.assigns.current_scope.user.id

    case Tasks.complete_task(task, user_id) do
      {:ok, task} ->
        {:noreply,
         socket
         |> stream_insert(:tasks, task)
         |> put_flash(:info, "Task completed")}

      {:error, _} ->
        {:noreply, put_flash(socket, :error, "Could not complete task")}
    end
  end

  def handle_event("cancel", %{"id" => id}, socket) do
    task = Tasks.get_task!(id)
    user_id = socket.assigns.current_scope.user.id

    case Tasks.cancel_task(task, user_id) do
      {:ok, task} ->
        {:noreply,
         socket
         |> stream_insert(:tasks, task)
         |> put_flash(:info, "Task cancelled")}

      {:error, _} ->
        {:noreply, put_flash(socket, :error, "Could not cancel task")}
    end
  end

  def handle_event("delete", %{"id" => id}, socket) do
    task = Tasks.get_task!(id)
    user_id = socket.assigns.current_scope.user.id

    case Tasks.delete_task(task, user_id) do
      {:ok, _} ->
        {:noreply,
         socket
         |> stream_delete(:tasks, task)
         |> put_flash(:info, "Task deleted")}

      {:error, _} ->
        {:noreply, put_flash(socket, :error, "Could not delete task")}
    end
  end

  @filter_options [
    {"All", "all"},
    {"Incomplete", "incomplete"},
    {"Completed", "completed"},
    {"Overdue", "overdue"},
    {"Today", "today"},
    {"This Week", "this_week"},
    {"This Month", "this_month"}
  ]

  @category_labels %{
    "general" => "General",
    "area" => "Area",
    "reservoir" => "Reservoir",
    "crop" => "Crop",
    "finance" => "Finance",
    "inventory" => "Inventory",
    "safety" => "Safety",
    "sanitation" => "Sanitation",
    "pestcontrol" => "Pest Control",
    "other" => "Other"
  }

  defp category_label(cat), do: Map.get(@category_labels, cat, cat)

  defp category_options do
    Enum.map(@category_labels, fn {key, label} -> {label, key} end)
    |> Enum.sort_by(&elem(&1, 0))
  end

  defp format_date(nil), do: "-"
  defp format_date(dt), do: Calendar.strftime(dt, "%b %d, %Y")

  defp priority_variant("urgent"), do: "danger"
  defp priority_variant(_), do: "default"

  defp status_variant("created"), do: "info"
  defp status_variant("completed"), do: "success"
  defp status_variant("cancelled"), do: "default"
  defp status_variant(_), do: "default"

  @impl true
  def render(assigns) do
    assigns = assign(assigns, :filter_options, @filter_options)

    ~H"""
    <Layouts.app flash={@flash} current_scope={@current_scope} active_page={:tasks}>
      <.header>
        {gettext("Tasks")}
        <:actions>
          <.button phx-click="new_task">
            <.icon name="hero-plus" class="size-4" /> {gettext("New Task")}
          </.button>
        </:actions>
      </.header>

      <div class="flex flex-col sm:flex-row gap-6 mt-6">
        <%!-- Filter sidebar --%>
        <nav class="sm:w-48 shrink-0">
          <ul class="space-y-1">
            <li :for={{label, value} <- @filter_options}>
              <button
                phx-click={JS.push("filter", value: %{filter: value})}
                class={[
                  "w-full text-left rounded-lg px-3 py-2 text-sm font-medium transition-colors cursor-pointer",
                  if(@filter == value,
                    do: "bg-primary text-white",
                    else: "text-text-light hover:bg-gray-100 hover:text-text"
                  )
                ]}
              >
                {label}
              </button>
            </li>
          </ul>
        </nav>

        <%!-- Task list --%>
        <div class="flex-1 min-w-0">
          <div :if={@streams.tasks |> Enum.any?()} class="space-y-3">
            <div
              :for={{dom_id, task} <- @streams.tasks}
              id={dom_id}
              class="rounded-xl border border-border bg-white p-4 shadow-sm hover:shadow-md transition-shadow"
            >
              <div class="flex items-start justify-between gap-3">
                <div class="flex-1 min-w-0">
                  <h3 class={[
                    "font-semibold text-text",
                    task.status == "completed" && "line-through opacity-60"
                  ]}>
                    {task.title}
                  </h3>
                  <p :if={task.description} class="mt-1 text-sm text-text-light line-clamp-2">
                    {task.description}
                  </p>
                  <div class="mt-2 flex flex-wrap gap-2">
                    <.badge variant={priority_variant(task.priority)}>
                      {String.capitalize(task.priority)}
                    </.badge>
                    <.badge variant={status_variant(task.status)}>
                      {String.capitalize(task.status)}
                    </.badge>
                    <.badge>{category_label(task.category)}</.badge>
                    <.badge :if={task.is_due} variant="danger">{gettext("Overdue")}</.badge>
                  </div>
                  <p :if={task.due_date} class="mt-2 text-xs text-text-light">
                    {gettext("Due")}: {format_date(task.due_date)}
                  </p>
                </div>
                <div class="flex gap-1 shrink-0">
                  <.button
                    :if={task.status == "created"}
                    variant="ghost"
                    phx-click={JS.push("complete", value: %{id: task.id})}
                    title={gettext("Complete")}
                  >
                    <.icon name="hero-check" class="size-4 text-success" />
                  </.button>
                  <.button
                    :if={task.status == "created"}
                    variant="ghost"
                    phx-click={JS.push("cancel", value: %{id: task.id})}
                    title={gettext("Cancel")}
                  >
                    <.icon name="hero-x-mark" class="size-4" />
                  </.button>
                  <.button
                    variant="ghost"
                    phx-click={JS.push("delete", value: %{id: task.id})}
                    data-confirm="Are you sure?"
                    title={gettext("Delete")}
                  >
                    <.icon name="hero-trash" class="size-4" />
                  </.button>
                </div>
              </div>
            </div>
          </div>

          <.empty_state
            :if={!(@streams.tasks |> Enum.any?())}
            icon="hero-clipboard-document-list"
            message={gettext("No tasks found.")}
          >
            <:actions>
              <.button phx-click="new_task">
                <.icon name="hero-plus" class="size-4" /> {gettext("New Task")}
              </.button>
            </:actions>
          </.empty_state>
        </div>
      </div>

      <.modal :if={@show_form} id="task-form-modal" show on_cancel={JS.push("close_form")}>
        <h2 class="text-lg font-semibold mb-4">{gettext("New Task")}</h2>
        <.form for={@form} phx-change="validate" phx-submit="save">
          <.input field={@form[:title]} label={gettext("Title")} required />
          <.input field={@form[:description]} type="textarea" label={gettext("Description")} />
          <.input
            field={@form[:priority]}
            type="select"
            label={gettext("Priority")}
            options={[{gettext("Normal"), "normal"}, {gettext("Urgent"), "urgent"}]}
          />
          <.input
            field={@form[:category]}
            type="select"
            label={gettext("Category")}
            options={category_options()}
          />
          <.input field={@form[:due_date]} type="datetime-local" label={gettext("Due Date")} />
          <div class="mt-4 flex justify-end gap-2">
            <.button type="button" variant="secondary" phx-click="close_form">
              {gettext("Cancel")}
            </.button>
            <.button type="submit">{gettext("Create")}</.button>
          </div>
        </.form>
      </.modal>
    </Layouts.app>
    """
  end
end
