import "../App.css";

function Index() {
  return (
    <>
      {/* this is needed so that app still has overflow */}
      <style>
        {`
          body{
            overflow:hidden;
          }
        `}
      </style>
      <span className="indexmain">
        <span className="mainitems" >
          <div>
            <h1>GYMJOURNAL</h1>
            <a href="http://localhost:8080/auth/google">Login with Google</a>
          </div>
          <div>
            <img src="https://www.ironmanmagazine.com/wp-content/uploads/Mike-Mentzers-Heavy-Duty-Workout-and-Diet-Plan.png" alt="picture of mike mentzer" />
          </div>
        </span>
      </span>
    </>
  );
}

export default Index;
