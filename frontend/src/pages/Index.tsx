import "../App.css";
import indexStyles from './index.module.css'
function Index() {
  return (
    <>
      <style>
        {`
          body{
            overflow:hidden;
          }
        `}
      </style>
      <span className={indexStyles['main']}>
          <div className={indexStyles['description']}>
            <h1>GYMJOURNAL</h1>
            <h4>A MINIMALISTIC WORKOUT LOGGING APP</h4>
            <a href="http://localhost:8080/auth/google">Login with Google</a>
          </div>
          <div>
            <img src="https://www.ironmanmagazine.com/wp-content/uploads/Mike-Mentzers-Heavy-Duty-Workout-and-Diet-Plan.png" alt="picture of mike mentzer" />
          </div>
      </span>
    </>
  );
}

export default Index;
