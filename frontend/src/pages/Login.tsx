const Login = () : JSX.Element => {
  return(
    <div className="auth-wrapper">
      <div className="card">
        <div className="row align-items-center text-center">
				  <div className="col-md-12">
            <div className="card-body">
              <h5 className="mb-3 f-w-400">Please, enter your credential.</h5>
              <div className="input-group mb-3">
                <span className="input-group-text"><i className="bi bi-person-circle"></i></span>
                <input type="text" className="form-control" placeholder="Username" />
              </div>
              <div className="input-group mb-4">
							  <span className="input-group-text"><i className="bi bi-lock"></i></span>
							  <input type="password" className="form-control" placeholder="Password" />
						  </div>
              <button className="btn btn-block btn-primary mb-4">Signin</button>
						  <p className="mb-0 text-muted">Don’t have an account? <a href="/" className="f-w-400">Signup</a></p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Login;
