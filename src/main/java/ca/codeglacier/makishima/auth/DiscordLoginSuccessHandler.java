package ca.codeglacier.makishima.auth;

import ca.codeglacier.makishima.model.DiscordUser;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClient;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientService;
import org.springframework.security.oauth2.client.authentication.OAuth2AuthenticationToken;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.security.web.authentication.AuthenticationSuccessHandler;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.sql.Date;
import java.util.NoSuchElementException;
import java.util.Optional;

@Component
public class DiscordLoginSuccessHandler implements AuthenticationSuccessHandler {

    private static final Logger logger = LoggerFactory.getLogger(DiscordLoginSuccessHandler.class);

    private final OAuth2AuthorizedClientService clientService;

    public DiscordLoginSuccessHandler(OAuth2AuthorizedClientService clientService) {
        this.clientService = clientService;
    }

    @Override
    public void onAuthenticationSuccess(HttpServletRequest request, HttpServletResponse response, Authentication authentication) throws IOException, ServletException {
        OAuth2AuthenticationToken token = (OAuth2AuthenticationToken) authentication;
        OAuth2User user = token.getPrincipal();
        OAuth2AuthorizedClient authorizedClient = clientService.loadAuthorizedClient(token.getAuthorizedClientRegistrationId(), user.getName());

        try {
            DiscordUser userRecord = new DiscordUser(
                    Optional.ofNullable(user.<Long>getAttribute("id")).orElseThrow(),
                    authorizedClient.getAccessToken().getTokenValue(),
                    Optional.ofNullable(authorizedClient.getRefreshToken()).orElseThrow().getTokenValue(),
                    new Date(0));
        } catch (NoSuchElementException e) {
            logger.warn(e.getMessage());
        }

        response.sendRedirect("/");
    }

}
