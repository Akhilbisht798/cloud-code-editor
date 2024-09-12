import { useEffect } from "react";
import { FC } from "react";
import { useAuth } from "@/state/userAuth";
import { useNavigate } from "react-router-dom";
import { SERVER } from "@/global";
import { Button } from "@/components/ui/button";

const Home: FC = () => {
  const navigate = useNavigate();
  const { email, setEmail } = useAuth();

  useEffect(() => {
    const url = SERVER + "/api/user";
    async function setUserAuth() {
      try {
        const res = await fetch(url, {
          method: "GET",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
        });
        const data = await res.json();
        if (res.ok) {
          setEmail(data.email);
        } else {
          navigate("/login");
        }
      } catch (err) {
        console.log(err);
        navigate("/login");
      }
    }
    setUserAuth();
  }, []);

  async function createNewProject(project: string) {
    try {
      const url = SERVER + "/api/createProject";
      const data = {
        userId: email,
        projectId: project,
      };
      const res = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });
      const resdata = await res.json();
      console.log(resdata);
    } catch (error) {
      console.log(error);
    }
  }

  if (email === "") {
    navigate("/login");
  }
  return (
    <>
      <div>This is home</div>
      <div>Welcome {email}</div>
      <Button onClick={createNewProject}>Create a New project</Button>
    </>
  );
};

export default Home;
