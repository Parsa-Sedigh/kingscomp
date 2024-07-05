We could write the services as interfaces and then mock the whole service.App and test it. Now the service.App can be used in
telegram or a CLI or http.

With this setup, we have a service.App and we can use it anywhere we want.