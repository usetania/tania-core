defmodule TaniaCoreWeb.Plugs.SetLocale do
  @moduledoc """
  Plug that sets the Gettext locale from the session or Accept-Language header.
  """

  import Plug.Conn

  @supported_locales ~w(en id)
  @default_locale "en"

  def init(opts), do: opts

  def call(conn, _opts) do
    locale =
      get_session(conn, "locale") ||
        get_preferred_locale(conn) ||
        @default_locale

    locale = if locale in @supported_locales, do: locale, else: @default_locale

    Gettext.put_locale(TaniaCoreWeb.Gettext, locale)
    conn |> put_session("locale", locale)
  end

  defp get_preferred_locale(conn) do
    case get_req_header(conn, "accept-language") do
      [header | _] -> parse_accept_language(header)
      _ -> nil
    end
  end

  defp parse_accept_language(header) do
    header
    |> String.split(",")
    |> Enum.map(fn part ->
      part
      |> String.trim()
      |> String.split(";")
      |> List.first()
      |> String.split("-")
      |> List.first()
      |> String.downcase()
    end)
    |> Enum.find(&(&1 in @supported_locales))
  end
end
