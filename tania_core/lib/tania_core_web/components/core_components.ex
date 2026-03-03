defmodule TaniaCoreWeb.CoreComponents do
  @moduledoc """
  Provides core UI components built with Tailwind CSS utilities.
  """
  use Phoenix.Component
  use Gettext, backend: TaniaCoreWeb.Gettext

  alias Phoenix.LiveView.JS

  @doc """
  Renders flash notices.

  ## Examples

      <.flash kind={:info} flash={@flash} />
      <.flash kind={:info} phx-mounted={show("#flash")}>Welcome Back!</.flash>
  """
  attr :id, :string, doc: "the optional id of flash container"
  attr :flash, :map, default: %{}, doc: "the map of flash messages to display"
  attr :title, :string, default: nil
  attr :kind, :atom, values: [:info, :error], doc: "used for styling and flash lookup"
  attr :rest, :global, doc: "the arbitrary HTML attributes to add to the flash container"

  slot :inner_block, doc: "the optional inner block that renders the flash message"

  def flash(assigns) do
    assigns = assign_new(assigns, :id, fn -> "flash-#{assigns.kind}" end)

    ~H"""
    <div
      :if={msg = render_slot(@inner_block) || Phoenix.Flash.get(@flash, @kind)}
      id={@id}
      phx-click={JS.push("lv:clear-flash", value: %{key: @kind}) |> hide("##{@id}")}
      role="alert"
      class="fixed top-4 right-4 z-50"
      {@rest}
    >
      <div class={[
        "w-80 sm:w-96 rounded-lg p-4 shadow-lg ring-1 flex items-start gap-3",
        @kind == :info && "bg-blue-50 text-blue-900 ring-blue-200",
        @kind == :error && "bg-red-50 text-red-900 ring-red-200"
      ]}>
        <.icon
          :if={@kind == :info}
          name="hero-information-circle"
          class="size-5 shrink-0 text-blue-500"
        />
        <.icon
          :if={@kind == :error}
          name="hero-exclamation-circle"
          class="size-5 shrink-0 text-red-500"
        />
        <div class="flex-1">
          <p :if={@title} class="font-semibold text-sm">{@title}</p>
          <p class="text-sm">{msg}</p>
        </div>
        <button type="button" class="group cursor-pointer" aria-label={gettext("close")}>
          <.icon name="hero-x-mark" class="size-5 opacity-40 group-hover:opacity-70" />
        </button>
      </div>
    </div>
    """
  end

  @doc """
  Renders a button.

  ## Examples

      <.button>Send!</.button>
      <.button phx-click="go" variant="primary">Send!</.button>
      <.button navigate={~p"/"}>Home</.button>
  """
  attr :rest, :global,
    include: ~w(href navigate patch method download name value disabled form type)

  attr :class, :string, default: nil
  attr :variant, :string, default: "primary", values: ~w(primary secondary danger ghost)
  slot :inner_block, required: true

  def button(%{rest: rest} = assigns) do
    variants = %{
      "primary" => "bg-primary text-white hover:bg-primary-dark focus:ring-primary/50",
      "secondary" =>
        "bg-white text-text border border-border hover:bg-gray-50 focus:ring-primary/50",
      "danger" => "bg-danger text-white hover:bg-danger-dark focus:ring-danger/50",
      "ghost" =>
        "bg-transparent text-text-light hover:text-text hover:bg-gray-100 focus:ring-primary/50"
    }

    base =
      "inline-flex items-center justify-center gap-2 rounded-lg px-4 py-2 text-sm font-semibold transition-colors focus:outline-none focus:ring-2 disabled:opacity-50 disabled:cursor-not-allowed phx-submit-loading:opacity-75 phx-click-loading:opacity-75"

    assigns =
      assign(assigns, :computed_class, [
        base,
        Map.fetch!(variants, assigns.variant),
        assigns.class
      ])

    if rest[:href] || rest[:navigate] || rest[:patch] do
      ~H"""
      <.link class={@computed_class} {@rest}>
        {render_slot(@inner_block)}
      </.link>
      """
    else
      ~H"""
      <button class={@computed_class} {@rest}>
        {render_slot(@inner_block)}
      </button>
      """
    end
  end

  @doc """
  Renders an input with label and error messages.

  A `Phoenix.HTML.FormField` may be passed as argument,
  which is used to retrieve the input name, id, and values.
  Otherwise all attributes may be passed explicitly.

  ## Types

  This function accepts all HTML input types, considering that:

    * You may also set `type="select"` to render a `<select>` tag

    * `type="checkbox"` is used exclusively to render boolean values

    * For live file uploads, see `Phoenix.Component.live_file_input/1`

  ## Examples

      <.input field={@form[:email]} type="email" />
      <.input name="my-input" errors={["oh no!"]} />
  """
  attr :id, :any, default: nil
  attr :name, :any
  attr :label, :string, default: nil
  attr :value, :any

  attr :type, :string,
    default: "text",
    values: ~w(checkbox color date datetime-local email file month number password
               search select tel text textarea time url week)

  attr :field, Phoenix.HTML.FormField,
    doc: "a form field struct retrieved from the form, for example: @form[:email]"

  attr :errors, :list, default: []
  attr :checked, :boolean, doc: "the checked flag for checkbox inputs"
  attr :prompt, :string, default: nil, doc: "the prompt for select inputs"
  attr :options, :list, doc: "the options to pass to Phoenix.HTML.Form.options_for_select/2"
  attr :multiple, :boolean, default: false, doc: "the multiple flag for select inputs"
  attr :class, :string, default: nil, doc: "the input class to use over defaults"
  attr :error_class, :string, default: nil, doc: "the input error class to use over defaults"

  attr :rest, :global,
    include: ~w(accept autocomplete capture cols disabled form list max maxlength min minlength
                multiple pattern placeholder readonly required rows size step)

  def input(%{field: %Phoenix.HTML.FormField{} = field} = assigns) do
    errors = if Phoenix.Component.used_input?(field), do: field.errors, else: []

    assigns
    |> assign(field: nil, id: assigns.id || field.id)
    |> assign(:errors, Enum.map(errors, &translate_error(&1)))
    |> assign_new(:name, fn -> if assigns.multiple, do: field.name <> "[]", else: field.name end)
    |> assign_new(:value, fn -> field.value end)
    |> input()
  end

  def input(%{type: "checkbox"} = assigns) do
    assigns =
      assign_new(assigns, :checked, fn ->
        Phoenix.HTML.Form.normalize_value("checkbox", assigns[:value])
      end)

    ~H"""
    <div class="mb-3">
      <label class="flex items-center gap-2 cursor-pointer">
        <input type="hidden" name={@name} value="false" disabled={@rest[:disabled]} />
        <input
          type="checkbox"
          id={@id}
          name={@name}
          value="true"
          checked={@checked}
          class={@class || "size-4 rounded border-border text-primary focus:ring-primary/50"}
          {@rest}
        />
        <span class="text-sm text-text">{@label}</span>
      </label>
      <.error :for={msg <- @errors}>{msg}</.error>
    </div>
    """
  end

  def input(%{type: "select"} = assigns) do
    ~H"""
    <div class="mb-3">
      <label>
        <span :if={@label} class="block text-sm font-medium text-text mb-1">{@label}</span>
        <select
          id={@id}
          name={@name}
          class={[
            @class ||
              "w-full rounded-lg border border-border bg-white px-3 py-2 text-sm text-text shadow-sm focus:border-primary focus:ring-1 focus:ring-primary/50 focus:outline-none",
            @errors != [] &&
              (@error_class || "border-danger focus:border-danger focus:ring-danger/50")
          ]}
          multiple={@multiple}
          {@rest}
        >
          <option :if={@prompt} value="">{@prompt}</option>
          {Phoenix.HTML.Form.options_for_select(@options, @value)}
        </select>
      </label>
      <.error :for={msg <- @errors}>{msg}</.error>
    </div>
    """
  end

  def input(%{type: "textarea"} = assigns) do
    ~H"""
    <div class="mb-3">
      <label>
        <span :if={@label} class="block text-sm font-medium text-text mb-1">{@label}</span>
        <textarea
          id={@id}
          name={@name}
          class={[
            @class ||
              "w-full rounded-lg border border-border bg-white px-3 py-2 text-sm text-text shadow-sm focus:border-primary focus:ring-1 focus:ring-primary/50 focus:outline-none",
            @errors != [] &&
              (@error_class || "border-danger focus:border-danger focus:ring-danger/50")
          ]}
          {@rest}
        >{Phoenix.HTML.Form.normalize_value("textarea", @value)}</textarea>
      </label>
      <.error :for={msg <- @errors}>{msg}</.error>
    </div>
    """
  end

  # All other inputs text, datetime-local, url, password, etc. are handled here...
  def input(assigns) do
    ~H"""
    <div class="mb-3">
      <label>
        <span :if={@label} class="block text-sm font-medium text-text mb-1">{@label}</span>
        <input
          type={@type}
          name={@name}
          id={@id}
          value={Phoenix.HTML.Form.normalize_value(@type, @value)}
          class={[
            @class ||
              "w-full rounded-lg border border-border bg-white px-3 py-2 text-sm text-text shadow-sm focus:border-primary focus:ring-1 focus:ring-primary/50 focus:outline-none",
            @errors != [] &&
              (@error_class || "border-danger focus:border-danger focus:ring-danger/50")
          ]}
          {@rest}
        />
      </label>
      <.error :for={msg <- @errors}>{msg}</.error>
    </div>
    """
  end

  # Helper used by inputs to generate form errors
  defp error(assigns) do
    ~H"""
    <p class="mt-1 flex items-center gap-1.5 text-sm text-danger">
      <.icon name="hero-exclamation-circle" class="size-4" />
      {render_slot(@inner_block)}
    </p>
    """
  end

  @doc """
  Renders a header with title.
  """
  slot :inner_block, required: true
  slot :subtitle
  slot :actions

  def header(assigns) do
    ~H"""
    <header class={[@actions != [] && "flex items-center justify-between gap-6", "pb-4"]}>
      <div>
        <h1 class="text-lg font-semibold leading-8 text-text">
          {render_slot(@inner_block)}
        </h1>
        <p :if={@subtitle != []} class="text-sm text-text-light">
          {render_slot(@subtitle)}
        </p>
      </div>
      <div class="flex-none">{render_slot(@actions)}</div>
    </header>
    """
  end

  @doc ~S"""
  Renders a modal dialog.

  ## Examples

      <.modal id="confirm-modal">
        Are you sure?
      </.modal>

  JS commands may be used to show/hide the modal:

      <button phx-click={show_modal("confirm-modal")}>Show</button>
  """
  attr :id, :string, required: true
  attr :show, :boolean, default: false
  attr :on_cancel, JS, default: %JS{}
  slot :inner_block, required: true

  def modal(assigns) do
    ~H"""
    <div
      id={@id}
      phx-mounted={@show && show_modal(@id)}
      phx-remove={hide_modal(@id)}
      data-cancel={JS.exec(@on_cancel, "phx-remove")}
      class="relative z-50 hidden"
    >
      <div id={"#{@id}-bg"} class="fixed inset-0 bg-black/50 transition-opacity" aria-hidden="true" />
      <div
        class="fixed inset-0 overflow-y-auto"
        aria-labelledby={"#{@id}-title"}
        aria-describedby={"#{@id}-description"}
        role="dialog"
        aria-modal="true"
        tabindex="0"
      >
        <div class="flex min-h-full items-center justify-center p-4">
          <div class="w-full max-w-lg">
            <.focus_wrap
              id={"#{@id}-container"}
              phx-window-keydown={JS.exec("data-cancel", to: "##{@id}")}
              phx-key="escape"
              phx-click-away={JS.exec("data-cancel", to: "##{@id}")}
              class="relative rounded-xl bg-white p-6 shadow-xl ring-1 ring-border"
            >
              <button
                phx-click={JS.exec("data-cancel", to: "##{@id}")}
                type="button"
                class="absolute top-4 right-4 cursor-pointer text-text-light hover:text-text"
                aria-label={gettext("close")}
              >
                <.icon name="hero-x-mark" class="size-5" />
              </button>
              <div id={"#{@id}-content"}>
                {render_slot(@inner_block)}
              </div>
            </.focus_wrap>
          </div>
        </div>
      </div>
    </div>
    """
  end

  @doc """
  Renders a card container.

  ## Examples

      <.card>
        <:header>Card Title</:header>
        Content here
      </.card>
  """
  attr :class, :string, default: nil
  attr :rest, :global
  slot :header
  slot :inner_block, required: true

  def card(assigns) do
    ~H"""
    <div class={["rounded-xl border border-border bg-white shadow-sm", @class]} {@rest}>
      <div :if={@header != []} class="border-b border-border px-5 py-3">
        <h3 class="text-base font-semibold text-text">{render_slot(@header)}</h3>
      </div>
      <div class="p-5">
        {render_slot(@inner_block)}
      </div>
    </div>
    """
  end

  @doc """
  Renders a badge.

  ## Examples

      <.badge variant="success">Active</.badge>
      <.badge variant="warning">Pending</.badge>
  """
  attr :variant, :string,
    default: "default",
    values: ~w(default primary success warning danger info)

  attr :class, :string, default: nil
  slot :inner_block, required: true

  def badge(assigns) do
    variants = %{
      "default" => "bg-gray-100 text-gray-700",
      "primary" => "bg-primary/10 text-primary",
      "success" => "bg-success/10 text-success-dark",
      "warning" => "bg-warning/10 text-warning-dark",
      "danger" => "bg-danger/10 text-danger",
      "info" => "bg-info/10 text-info-dark"
    }

    assigns = assign(assigns, :variant_class, Map.fetch!(variants, assigns.variant))

    ~H"""
    <span class={[
      "inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium",
      @variant_class,
      @class
    ]}>
      {render_slot(@inner_block)}
    </span>
    """
  end

  @doc """
  Renders tabs.

  ## Examples

      <.tabs tabs={[
        %{label: "Active", value: "active", active: true},
        %{label: "Archived", value: "archived", active: false}
      ]} on_click="switch_tab" />
  """
  attr :tabs, :list, required: true, doc: "list of %{label, value, active} maps"
  attr :on_click, :string, required: true, doc: "event name to push on tab click"

  def tabs(assigns) do
    ~H"""
    <div class="border-b border-border">
      <nav class="-mb-px flex gap-4" aria-label="Tabs">
        <button
          :for={tab <- @tabs}
          phx-click={@on_click}
          phx-value-tab={tab.value}
          class={[
            "whitespace-nowrap border-b-2 py-3 px-1 text-sm font-medium transition-colors cursor-pointer",
            if(tab[:active],
              do: "border-primary text-primary",
              else: "border-transparent text-text-light hover:border-border hover:text-text"
            )
          ]}
        >
          {tab.label}
        </button>
      </nav>
    </div>
    """
  end

  @doc """
  Renders an empty state placeholder.

  ## Examples

      <.empty_state icon="hero-rectangle-stack" message="No crops found" />
  """
  attr :icon, :string, default: "hero-inbox"
  attr :message, :string, required: true
  attr :class, :string, default: nil
  slot :actions

  def empty_state(assigns) do
    ~H"""
    <div class={["flex flex-col items-center justify-center py-12 text-center", @class]}>
      <.icon name={@icon} class="size-12 text-text-light/50 mb-4" />
      <p class="text-sm text-text-light">{@message}</p>
      <div :if={@actions != []} class="mt-4">
        {render_slot(@actions)}
      </div>
    </div>
    """
  end

  @doc """
  Renders a sidebar navigation link.

  ## Examples

      <.sidebar_link icon="hero-home" label="Dashboard" href="/" active={@active_page == :dashboard} />
  """
  attr :icon, :string, required: true
  attr :label, :string, required: true
  attr :active, :boolean, default: false
  attr :rest, :global, include: ~w(href navigate patch)

  def sidebar_link(assigns) do
    ~H"""
    <.link
      class={[
        "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
        if(@active,
          do: "bg-sidebar-active text-white",
          else: "text-white/70 hover:bg-sidebar-hover hover:text-white"
        )
      ]}
      {@rest}
    >
      <.icon name={@icon} class="size-5" />
      <span>{@label}</span>
    </.link>
    """
  end

  @doc """
  Renders a statistics card for the dashboard.

  ## Examples

      <.stat_card label="Active Crops" value={42} icon="hero-squares-2x2" />
  """
  attr :label, :string, required: true
  attr :value, :any, required: true
  attr :icon, :string, default: nil
  attr :class, :string, default: nil

  def stat_card(assigns) do
    ~H"""
    <div class={["rounded-xl border border-border bg-white p-5 shadow-sm", @class]}>
      <div class="flex items-center justify-between">
        <div>
          <p class="text-sm text-text-light">{@label}</p>
          <p class="mt-1 text-2xl font-bold text-text">{@value}</p>
        </div>
        <div :if={@icon} class="rounded-lg bg-primary/10 p-3">
          <.icon name={@icon} class="size-6 text-primary" />
        </div>
      </div>
    </div>
    """
  end

  @doc """
  Renders pagination controls.

  ## Examples

      <.pagination page={@page} total_pages={@total_pages} />
  """
  attr :page, :integer, required: true
  attr :total_pages, :integer, required: true
  attr :class, :string, default: nil

  def pagination(assigns) do
    ~H"""
    <nav
      :if={@total_pages > 1}
      class={["flex items-center justify-center gap-1", @class]}
      aria-label="Pagination"
    >
      <button
        phx-click="paginate"
        phx-value-page={@page - 1}
        disabled={@page <= 1}
        class="rounded-lg p-2 text-text-light hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <.icon name="hero-chevron-left" class="size-4" />
      </button>
      <span class="px-3 py-1 text-sm text-text">
        {gettext("Page %{page} of %{total}", page: @page, total: @total_pages)}
      </span>
      <button
        phx-click="paginate"
        phx-value-page={@page + 1}
        disabled={@page >= @total_pages}
        class="rounded-lg p-2 text-text-light hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <.icon name="hero-chevron-right" class="size-4" />
      </button>
    </nav>
    """
  end

  @doc """
  Renders a timeline for activity display.

  ## Examples

      <.timeline>
        <.timeline_item icon="hero-arrow-right" time={~U[2024-01-01 12:00:00Z]}>
          Crop moved to Area B
        </.timeline_item>
      </.timeline>
  """
  slot :inner_block, required: true

  def timeline(assigns) do
    ~H"""
    <div class="flow-root">
      <ul role="list" class="-mb-8">
        {render_slot(@inner_block)}
      </ul>
    </div>
    """
  end

  @doc """
  Renders a single timeline item.
  """
  attr :icon, :string, default: "hero-clock"
  attr :time, :any, default: nil
  attr :last, :boolean, default: false
  slot :inner_block, required: true

  def timeline_item(assigns) do
    ~H"""
    <li>
      <div class="relative pb-8">
        <span
          :if={!@last}
          class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-border"
          aria-hidden="true"
        />
        <div class="relative flex gap-3">
          <div class="flex size-8 shrink-0 items-center justify-center rounded-full bg-primary/10 ring-4 ring-white">
            <.icon name={@icon} class="size-4 text-primary" />
          </div>
          <div class="flex-1 min-w-0 pt-0.5">
            <p class="text-sm text-text">{render_slot(@inner_block)}</p>
            <p :if={@time} class="mt-0.5 text-xs text-text-light">{@time}</p>
          </div>
        </div>
      </div>
    </li>
    """
  end

  @doc """
  Renders a table with generic styling.

  ## Examples

      <.table id="users" rows={@users}>
        <:col :let={user} label="id">{user.id}</:col>
        <:col :let={user} label="username">{user.username}</:col>
      </.table>
  """
  attr :id, :string, required: true
  attr :rows, :list, required: true
  attr :row_id, :any, default: nil, doc: "the function for generating the row id"
  attr :row_click, :any, default: nil, doc: "the function for handling phx-click on each row"

  attr :row_item, :any,
    default: &Function.identity/1,
    doc: "the function for mapping each row before calling the :col and :action slots"

  slot :col, required: true do
    attr :label, :string
  end

  slot :action, doc: "the slot for showing user actions in the last table column"

  def table(assigns) do
    assigns =
      with %{rows: %Phoenix.LiveView.LiveStream{}} <- assigns do
        assign(assigns, row_id: assigns.row_id || fn {id, _item} -> id end)
      end

    ~H"""
    <div class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead class="border-b border-border">
          <tr>
            <th
              :for={col <- @col}
              class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-text-light"
            >
              {col[:label]}
            </th>
            <th :if={@action != []} class="px-4 py-3">
              <span class="sr-only">{gettext("Actions")}</span>
            </th>
          </tr>
        </thead>
        <tbody
          id={@id}
          phx-update={is_struct(@rows, Phoenix.LiveView.LiveStream) && "stream"}
          class="divide-y divide-border"
        >
          <tr
            :for={row <- @rows}
            id={@row_id && @row_id.(row)}
            class="hover:bg-gray-50 transition-colors"
          >
            <td
              :for={col <- @col}
              phx-click={@row_click && @row_click.(row)}
              class={["px-4 py-3 text-text", @row_click && "cursor-pointer"]}
            >
              {render_slot(col, @row_item.(row))}
            </td>
            <td :if={@action != []} class="px-4 py-3 w-0 font-semibold">
              <div class="flex gap-4">
                <%= for action <- @action do %>
                  {render_slot(action, @row_item.(row))}
                <% end %>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    """
  end

  @doc """
  Renders a data list.

  ## Examples

      <.list>
        <:item title="Title">{@post.title}</:item>
        <:item title="Views">{@post.views}</:item>
      </.list>
  """
  slot :item, required: true do
    attr :title, :string, required: true
  end

  def list(assigns) do
    ~H"""
    <dl class="divide-y divide-border">
      <div :for={item <- @item} class="flex gap-4 py-3">
        <dt class="w-1/4 flex-none text-sm font-medium text-text-light">{item.title}</dt>
        <dd class="text-sm text-text">{render_slot(item)}</dd>
      </div>
    </dl>
    """
  end

  @doc """
  Renders a [Heroicon](https://heroicons.com).

  Heroicons come in three styles – outline, solid, and mini.
  By default, the outline style is used, but solid and mini may
  be applied by using the `-solid` and `-mini` suffix.

  ## Examples

      <.icon name="hero-x-mark" />
      <.icon name="hero-arrow-path" class="ml-1 size-3 motion-safe:animate-spin" />
  """
  attr :name, :string, required: true
  attr :class, :string, default: "size-4"

  def icon(%{name: "hero-" <> _} = assigns) do
    ~H"""
    <span class={[@name, @class]} />
    """
  end

  ## JS Commands

  def show(js \\ %JS{}, selector) do
    JS.show(js,
      to: selector,
      time: 300,
      transition:
        {"transition-all ease-out duration-300",
         "opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95",
         "opacity-100 translate-y-0 sm:scale-100"}
    )
  end

  def hide(js \\ %JS{}, selector) do
    JS.hide(js,
      to: selector,
      time: 200,
      transition:
        {"transition-all ease-in duration-200", "opacity-100 translate-y-0 sm:scale-100",
         "opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"}
    )
  end

  def show_modal(js \\ %JS{}, id) when is_binary(id) do
    js
    |> JS.show(to: "##{id}")
    |> JS.show(
      to: "##{id}-bg",
      time: 300,
      transition: {"transition-all ease-out duration-300", "opacity-0", "opacity-100"}
    )
    |> show("##{id}-container")
    |> JS.add_class("overflow-hidden", to: "body")
    |> JS.focus_first(to: "##{id}-content")
  end

  def hide_modal(js \\ %JS{}, id) do
    js
    |> JS.hide(
      to: "##{id}-bg",
      transition: {"transition-all ease-in duration-200", "opacity-100", "opacity-0"}
    )
    |> hide("##{id}-container")
    |> JS.hide(to: "##{id}", transition: {"block", "block", "hidden"})
    |> JS.remove_class("overflow-hidden", to: "body")
    |> JS.pop_focus()
  end

  @doc """
  Translates an error message using gettext.
  """
  def translate_error({msg, opts}) do
    if count = opts[:count] do
      Gettext.dngettext(TaniaCoreWeb.Gettext, "errors", msg, msg, count, opts)
    else
      Gettext.dgettext(TaniaCoreWeb.Gettext, "errors", msg, opts)
    end
  end

  @doc """
  Translates the errors for a field from a keyword list of errors.
  """
  def translate_errors(errors, field) when is_list(errors) do
    for {^field, {msg, opts}} <- errors, do: translate_error({msg, opts})
  end
end
