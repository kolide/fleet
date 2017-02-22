## Kolide 1.0.1 (February 22, 2016) ##
*   Improve platform detection accuracy

    Previously we were using `build_platform`, which does not always properly
    reflect the platform of the host running osquery.

*   Fix bugs where query links in the pack sidebar pointed to the wrong queries

*   Improve MySQL compatibility with stricter configurations

    In some MySQL configurations, using a GROUP BY that doesn't refer to every
    column in the SELECT will throw errors. Replace the use of GROUP BY with SELECT
    DISTINCT as this is also more clear as to the intentions of the query.

*   Allow users to edit the name and description of host labels

*   Add basic table autocompletion when typing in the query composer.

*   Support MySQL client certificate authentication

    More details can be found in the [Configuring the Kolide binary docs](https://docs.kolide.co/kolide/1.0.1/infrastructure/configuring-the-kolide-binary.html)

*   Improve security for user-initiated email address changes

    This improvement ensures that only users who own an email address and are
    logged in as the user who initiated the change can confirm the new email.

    Previously it was possible for Administrators to also confirm these changes
    by clicking the confirmation link.

*   Fix an issue where the setup form rejects passwords with certain characters

    This change resolves an issue where certain special characters like "."
    where rejected by the client-side JS that controls the setup form.

*   Automatically login the user once initial setup is completed
