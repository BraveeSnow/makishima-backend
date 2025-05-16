package ca.codeglacier.makishima.auth;

import ca.codeglacier.makishima.entity.DiscordUser;
import ca.codeglacier.makishima.service.DiscordUserService;
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
import java.util.Date;
import java.util.NoSuchElementException;
import java.util.Optional;

@Component
public class DiscordLoginSuccessHandler implements AuthenticationSuccessHandler {

    private static final Logger logger = LoggerFactory.getLogger(DiscordLoginSuccessHandler.class);

    private final OAuth2AuthorizedClientService clientService;
    private final DiscordUserService discordUserService;

    public DiscordLoginSuccessHandler(OAuth2AuthorizedClientService clientService, DiscordUserService discordUserService) {
        this.clientService = clientService;
        this.discordUserService = discordUserService;
    }

    @Override
    public void onAuthenticationSuccess(HttpServletRequest request, HttpServletResponse response, Authentication authentication) throws IOException {
        OAuth2AuthenticationToken token = (OAuth2AuthenticationToken) authentication;
        OAuth2User user = token.getPrincipal();
        OAuth2AuthorizedClient authorizedClient = clientService.loadAuthorizedClient(token.getAuthorizedClientRegistrationId(), user.getName());

        try {
            discordUserService.save(new DiscordUser(
                    Long.parseLong(Optional.ofNullable(user.getAttribute("id")).orElseThrow().toString()),
                    authorizedClient.getAccessToken().getTokenValue(),
                    Optional.ofNullable(authorizedClient.getRefreshToken()).orElseThrow().getTokenValue(),
                    Date.from(Optional.ofNullable(authorizedClient.getAccessToken().getExpiresAt()).orElseThrow())));
        } catch (NoSuchElementException e) {
            logger.warn(e.getMessage());
        }

        response.sendRedirect("/");
    }

}
