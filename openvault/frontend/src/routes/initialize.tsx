import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardTitle
} from '@/components/ui/card';
import {
  createFileRoute,
  useNavigate,
  useRouter
} from '@tanstack/react-router';
import { z } from 'zod';
import { useForm } from '@tanstack/react-form';
import { toast } from 'sonner';
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel
} from '@/components/ui/field';
import { Input } from '@/components/ui/input';
import { Separator } from '@/components/ui/separator';
import { IdCardLanyard, Loader } from 'lucide-react';
import { InitOptions } from '@openvault/openvault/internal/fs';
import { useInitialize } from '@/context/initialize';

export const Route = createFileRoute('/initialize')({
  component: InitializeForm
});

const formSchema = z
  .object({
    email: z.email(),
    firstName: z.string().min(1, 'First name is required'),
    lastName: z.string().min(1, 'Last name is required'),
    password: z.string().min(8, 'Password must be at least 8 characters long'),
    confirmPassword: z
      .string()
      .min(8, 'Confirm Password must be at least 8 characters long')
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword']
  });

function InitializeForm() {
  const navigate = useNavigate();
  const router = useRouter();
  const { initialize } = useInitialize();
  const form = useForm({
    defaultValues: {
      email: '',
      firstName: '',
      lastName: '',
      password: '',
      confirmPassword: ''
    },
    validators: {
      onSubmit: formSchema,
      onBlur: formSchema,
      onChange: formSchema
    },
    onSubmit: async ({ value }) => {
      try {
        const data = new InitOptions();
        data.Email = value.email;
        data.FirstName = value.firstName;
        data.LastName = value.lastName;
        data.Password = value.password;
        await initialize(data);
        toast.success('Account created successfully!');
        await router.invalidate();
        navigate({
          to: '/lock',
          search: {
            redirect: '/all'
          }
        });
      } catch (error) {
        toast.error('Failed to create account', {
          description: (
            <code className="font-mono break-all">{String(error)}</code>
          )
        });
      }
    }
  });

  return (
    <div className="flex h-full w-full items-center justify-center">
      <Card className="w-1/2">
        <CardContent className="flex flex-col gap-y-4">
          <CardTitle className="text-xl whitespace-nowrap">
            Welcome to OpenVault!
          </CardTitle>
          <CardDescription>
            Let&apos;s get you started with your account. Don&apos;t worry, this
            info is only used to help you identify and secure your account.
          </CardDescription>
          <Separator />
          <form
            id="initialize-form"
            onSubmit={(e) => {
              e.preventDefault();
              form.handleSubmit();
            }}
          >
            <FieldGroup className="gap-4">
              <div className="grid grid-cols-2 gap-4">
                <form.Field
                  name="firstName"
                  children={(field) => {
                    const isInvalid =
                      field.state.meta.isTouched && !field.state.meta.isValid;
                    return (
                      <Field data-invalid={isInvalid}>
                        <FieldLabel htmlFor={field.name}>First Name</FieldLabel>
                        <Input
                          id={field.name}
                          value={field.state.value}
                          onBlur={field.handleBlur}
                          onChange={(e) => field.handleChange(e.target.value)}
                          aria-invalid={isInvalid}
                          placeholder="Joe"
                          autoComplete="off"
                        />
                        {isInvalid && (
                          <FieldError errors={field.state.meta.errors} />
                        )}
                      </Field>
                    );
                  }}
                />
                <form.Field
                  name="lastName"
                  children={(field) => {
                    const isInvalid =
                      field.state.meta.isTouched && !field.state.meta.isValid;
                    return (
                      <Field data-invalid={isInvalid}>
                        <FieldLabel htmlFor={field.name}>Last Name</FieldLabel>
                        <Input
                          id={field.name}
                          value={field.state.value}
                          onBlur={field.handleBlur}
                          onChange={(e) => field.handleChange(e.target.value)}
                          aria-invalid={isInvalid}
                          placeholder="Shmoe"
                          autoComplete="off"
                        />
                        {isInvalid && (
                          <FieldError errors={field.state.meta.errors} />
                        )}
                      </Field>
                    );
                  }}
                />
              </div>
              <form.Field
                name="email"
                children={(field) => {
                  const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid;
                  return (
                    <Field data-invalid={isInvalid}>
                      <FieldLabel htmlFor={field.name}>Email</FieldLabel>
                      <Input
                        id={field.name}
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        aria-invalid={isInvalid}
                        placeholder="jshmoe@example.com"
                        autoComplete="off"
                      />
                      {isInvalid && (
                        <FieldError errors={field.state.meta.errors} />
                      )}
                    </Field>
                  );
                }}
              />
              <form.Field
                name="password"
                children={(field) => {
                  const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid;
                  return (
                    <Field data-invalid={isInvalid}>
                      <FieldLabel htmlFor={field.name}>Password</FieldLabel>
                      <Input
                        id={field.name}
                        type="password"
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        aria-invalid={isInvalid}
                        placeholder="****************"
                        autoComplete="off"
                      />
                      {isInvalid && (
                        <FieldError errors={field.state.meta.errors} />
                      )}
                    </Field>
                  );
                }}
              />
              <form.Field
                name="confirmPassword"
                children={(field) => {
                  const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid;
                  return (
                    <Field data-invalid={isInvalid}>
                      <FieldLabel htmlFor={field.name}>
                        Confirm Password
                      </FieldLabel>
                      <Input
                        id={field.name}
                        type="password"
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        aria-invalid={isInvalid}
                        placeholder="****************"
                        autoComplete="off"
                      />
                      {isInvalid && (
                        <FieldError errors={field.state.meta.errors} />
                      )}
                    </Field>
                  );
                }}
              />
              <form.Subscribe
                selector={(state) => [
                  state.canSubmit,
                  state.isPristine,
                  state.isSubmitting
                ]}
                children={([canSubmit, isPristine, isSubmitting]) => (
                  <Button
                    className="w-full"
                    disabled={!canSubmit || isPristine || isSubmitting}
                  >
                    {isSubmitting ? (
                      <Loader className="size-4 animate-spin" />
                    ) : (
                      <IdCardLanyard className="size-4" />
                    )}{' '}
                    Create Account
                  </Button>
                )}
              />
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
