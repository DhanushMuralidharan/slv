---
sidebar_position: 1
---

# Environment

## What is an Environment in SLV?
## Environment

In SLV, an **environment** is treated as a unique **identity**.

Each environment is associated with an **asymmetric key pair**:

- The **secret key** is securely held and used to decrypt secrets.
- The **public key**, conceptually, allows others to encrypt secrets in a way that only that environment can decrypt.

This is a simplified way to understand how environments function in SLV. A deeper understanding emerges through related concepts such as vaults, profiles, and access flows.

At its core, an environment acts as a **cryptographic identity** that can **decrypt** secrets that are specifically intended for that identity. However, an environment is more than just a cryptographic key. As the name suggests, it represents the **context in which an application is run**, with the necessary secrets loaded and ready for use. It is the runtime surface — whether it's a developer laptop, CI pipeline, container, or server — where those secrets are finally consumed.

---

## Components of an Environment
| **Component** | **Description** |
|-|-|
| **Public Key** | Used to share a secret with the environment. Starts with `SLV_EPK`|
| **Name** | Name of the Environment. |
| **Email** | Email Address Associated with the environment. |
| **Type** | Type of the environment (More about this later). |
| **Tags** | Used for adding additional metadata or labels for the environment. |
| **EDS (Environment Definition String)** | A string that is generated by SLV with all the above information encoded in it. This can be used to share your environment info in a standardised way. Starts with `SLV_EDS`|
| **Secret** | This can be a **secret key** (`SLV_ESK`) or a **secret binding** (`SLV_ESB`) which corresponds to the secret key component that is used to decrypt secrets. (Losing this can make the environment useless) |

---

## Types of Environment

### User Environment

The **user environment** is intended for cases where SLV is accessed from a laptop, desktop, or other personal end-user device.

In this setup, SLV uses a **password provided by the user** to derive a cryptographic token known as the **secret binding**. This secret binding serves as the identity of the environment and is used to decrypt secrets associated with it.

- The password is used once during setup to generate the secret binding.
- The generated **secret binding must be stored securely** by the user.
- SLV does **not store the password** or allow re-generation of the secret binding later.

The secret binding is essential for accessing secrets intended for that environment. If lost, the binding cannot be recovered.
This approach ensures that secrets can be securely accessed from personal devices without central storage of credentials.

### Service Environment
The **service environment** is used when SLV is accessed from automated or non-interactive systems (Typically in production) such as GitHub Actions, CI/CD pipelines, or Kubernetes clusters.

In this case, SLV generates a **secret key** directly, which acts as the identity for the environment.

- This secret key is used to decrypt secrets during runtime.
- It must be **stored securely**, as it is critical to accessing the vault.
- SLV does **not retain or regenerate** the secret key once issued.

While managing the raw secret key works, it may not always be convenient or secure in automation-heavy workflows. To address this, SLV supports integration with external **Key Management Services (KMS)** such as AWS KMS, Google Cloud KMS, or others.

By leveraging KMS:
- The secret key can be binded by the cloud provider.
- SLV can retrieve the secret binding and unbind the key at runtime, reducing the need for manual secret management.

This setup makes the service environment both secure and automation-friendly — ideal for integrating SLV into production pipelines and infrastructure.

### Root Environment

The **root environment** is a special type of environment that can be configured for a profile (details on profiles are covered later).

When a root environment is set up, SLV automatically shares **all newly created vaults** with this environment by default. This behavior acts as a safeguard to prevent accidental loss of access to secrets.

In the unlikely event that an environment’s **secret key** or **secret binding** is lost or becomes unrecoverable, the root environment still retains access to the vault and its secrets.

This ensures that there is always a fallback environment with full access, making the system more resilient and easier to recover from critical failures.



