package ca.codeglacier.makishima.service;

import ca.codeglacier.makishima.entity.DiscordUser;
import ca.codeglacier.makishima.repository.DiscordUserRepository;
import org.springframework.stereotype.Service;

@Service
public class DiscordUserService {

    private final DiscordUserRepository userRepository;

    public DiscordUserService(DiscordUserRepository userRepository) {
        this.userRepository = userRepository;
    }

    public void save(DiscordUser user) {
        userRepository.save(user);
    }

}
