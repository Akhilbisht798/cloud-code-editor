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
import { SERVER } from "@/global";
import { Link, useNavigate } from "react-router-dom";
import { useState } from "react";

const RegisterSchema = z
  .object({
    email: z
      .string()
      .min(1, { message: "This field has to be filled." })
      .email("this is not a valid email"),
    // .refine((e) => e === "abcd@fg.com", "This email is not in our database"),
    password: z.string().min(5),
    confirmPassword: z.string().min(5),
  })
  .superRefine((data, ctx) => {
    if (data.confirmPassword != data.password) {
      ctx.addIssue({
        code: "custom",
        message: "The password did not match",
        path: ["confirmPassword"],
      });
    }
  });

const Register: FC = () => {
  const [regiesterd, setRegistered] = useState(false);
  const navigate = useNavigate();
  const form = useForm<z.infer<typeof RegisterSchema>>({
    resolver: zodResolver(RegisterSchema),
    defaultValues: {
      email: "",
      password: "",
      confirmPassword: "",
    },
  });

  function onSubmit(values: z.infer<typeof RegisterSchema>) {
    const url = SERVER + "/api/register";
    const data = {
      email: values.email,
      password: values.password,
    };
    fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    })
      .then((res) => res.json())
      .then((res) => {
        console.log(res);
        setRegistered(true);
      })
      .catch((error) => {
        console.log("Error: ", error);
      });
  }

  if (regiesterd) {
    console.log("redirect to login");
    navigate("/login");
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-[#1a1b1e] dark:bg-[#1a1b1e] text-white">
      <div className="w-full max-w-md p-6 bg-[#2c2d31] rounded-lg shadow-lg dark:bg-[#2c2d31]">
        <div className="text-center">
          <h1 className="text-3xl font-bold dark:text-primary-foreground">
            Register Now
          </h1>
          <p className="text-muted-foreground">Create a new account</p>
        </div>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-white">Email</FormLabel>
                  <FormControl>
                    <Input placeholder="abc@mail.com" {...field} />
                  </FormControl>
                  <FormDescription>Enter a valid email</FormDescription>
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
                    <Input placeholder="Enter a Strong Password" {...field} />
                  </FormControl>
                  <FormDescription>Enter a strong password</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="confirmPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Confirm Password</FormLabel>
                  <FormControl>
                    <Input placeholder="Re-Enter a your Password" {...field} />
                  </FormControl>
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
        <div className="text-center text-muted-foreground text-white">
          <p>
            Already have a account?{" "}
            <Link to="/login" className="font-medium hover:underline">
              Login
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
};

export default Register;
