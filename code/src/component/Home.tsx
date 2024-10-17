import { useEffect } from "react";
import { FC } from "react";
import { useAuth } from "@/state/userAuth";
import { useNavigate } from "react-router-dom";
import { SERVER } from "@/global";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useCurrentProject } from "@/state/currentProject";

const Home: FC = () => {
  const navigate = useNavigate();
  const { email, setEmail } = useAuth();
  const { rootDir, setRootDir } = useCurrentProject();

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

  function setProjectHandler(e: any) {
    setRootDir(e.target.value);
  }

  async function createNewProject() {
    try {
      const url = SERVER + "/api/start";
      console.log("Project: ", rootDir);
      const data = {
        userId: email,
        projectId: rootDir,
      };
      const res = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
        credentials: "include",
      });
      const resdata = await res.json();
      console.log("ip: " + resdata);
      navigate("/project");
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
      <Input
        type="text"
        onChange={(e) => {
          setProjectHandler(e);
        }}
      />
      <Button onClick={createNewProject}>Create a New project</Button>
    </>
  );
};

export default Home;
