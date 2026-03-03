defmodule TaniaCore.Repo do
  use Ecto.Repo,
    otp_app: :tania_core,
    adapter: Ecto.Adapters.SQLite3
end
