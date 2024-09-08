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
  const form = useForm<z.infer<typeof RegisterSchema>>({
    resolver: zodResolver(RegisterSchema),
    defaultValues: {
      email: "",
      password: "",
      confirmPassword: "",
    },
  });

  function onSubmit(values: z.infer<typeof RegisterSchema>) {
    console.log(values);
  }

  return (
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

        <Button type="submit">Submit</Button>
      </form>
    </Form>
  );
};

export default Register;
