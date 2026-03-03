defmodule TaniaCoreWeb.UserSessionHTML do
  use TaniaCoreWeb, :html

  embed_templates "user_session_html/*"

  defp local_mail_adapter? do
    Application.get_env(:tania_core, TaniaCore.Mailer)[:adapter] == Swoosh.Adapters.Local
  end
end
