# Script for populating the database. You can run it as:
#
#     mix run priv/repo/seeds.exs

alias TaniaCore.Accounts

# Create default admin user
case Accounts.get_user_by_email("admin@tania.local") do
  nil ->
    {:ok, user} = Accounts.register_user(%{email: "admin@tania.local"})

    # Set password directly
    user
    |> Ecto.Changeset.change(%{
      hashed_password: Bcrypt.hash_pwd_salt("tania12345678"),
      confirmed_at: DateTime.utc_now() |> DateTime.truncate(:second)
    })
    |> TaniaCore.Repo.update!()

    IO.puts("Created default user: admin@tania.local / tania12345678")

  _user ->
    IO.puts("Default user admin@tania.local already exists")
end
