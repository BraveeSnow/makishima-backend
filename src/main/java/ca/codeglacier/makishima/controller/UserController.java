package ca.codeglacier.makishima.controller;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;
import java.util.Objects;

@RestController
public class UserController {

    private static final String DISCORD_AVATAR_URL = "https://cdn.discordapp.com/avatars/%s/%s.webp";

    @GetMapping("/user")
    public ResponseEntity<Map<String, String>> getUser(@AuthenticationPrincipal OAuth2User user) {
        if (user == null) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).build();
        }

        return ResponseEntity.ok(Map.of(
                "username", Objects.requireNonNull(user.getAttribute("username")),
                // Craft avatar URL according to Discord documentation
                // https://discord.com/developers/docs/reference#image-formatting
                "avatarUrl", DISCORD_AVATAR_URL.formatted(user.getAttribute("id"), Objects.requireNonNull(user.getAttribute("avatar")))
        ));
    }

}
