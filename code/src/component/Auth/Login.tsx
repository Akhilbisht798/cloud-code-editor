import { FC } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Link, useNavigate } from "react-router-dom";
import { SERVER } from "@/global";
import { useState } from "react";

const LoginSchema = z.object({
  email: z
    .string()
    .min(1, { message: "This field has to be filled." })
    .email("this is not a valid email"),
  // .refine((e) => e === "abcd@fg.com", "This email is not in our database"),
  password: z.string().min(5),
});

const Login: FC = () => {
  const [isError, setIsError] = useState("");
  const navigate = useNavigate();
  const form = useForm<z.infer<typeof LoginSchema>>({
    resolver: zodResolver(LoginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  async function onSubmit(values: z.infer<typeof LoginSchema>) {
    const url = SERVER + "/api/login";
    const data = {
      email: values.email,
      password: values.password,
    };
    try {
      const res = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(data),
      });
      if (res.ok) {
        // const body = await res.json();
        // const jwt = body['jwt'];
        navigate("/");
      } else {
        const errorText = await res.text();
        setIsError(errorText);
      }
    } catch (err) {
      console.log("error: ", err);
    }
  }

  return (
    <div className="flex items-center justify-center h-screen  bg-[#1a1b1e] dark:bg-[#1a1b1e] text-white font-sans">
      <div className="w-full max-w-md p-6 bg-[#2c2d31] rounded-lg shadow-lg dark:bg-[#2c2d31]">
        <div className="text-center">
          <h1 className="text-3xl font-bold dark:text-primary-foreground">
            Welcome Back
          </h1>
          <p className="text-muted-foreground">Sign in to your account</p>
        </div>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input placeholder="abc@mail.com" {...field} />
                  </FormControl>
                  <FormDescription>Enter a your email</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>password</FormLabel>
                  <FormControl>
                    <Input placeholder="Enter a your Password" {...field} />
                  </FormControl>
                  <FormDescription>Enter a your password</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button
              type="submit"
              className="w-full bg-primary text-primary-foreground hover:bg-primary/90 dark:bg-primary-foreground dark:text-primary hover:dark:bg-primary/90"
            >
              Submit
            </Button>
          </form>
        </Form>
        {isError !== "" && <span>{isError}</span>}
        <div className="text-center text-muted-foreground text-white">
          <p>
            Don't have an account?{" "}
            <Link to="/register" className="font-medium hover:underline">
              Register
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
};

export default Login;
